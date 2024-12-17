package main

import (
	"bytes"
	"encoding/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// A WikiConfig is a (de)serializable config which holds extensible settings.
type WikiConfig struct {
	Title     string `json:"title"`
	Surface   string `json:"theme-background"`
	Surface2  string `json:"theme-background2"`
	Surface3  string `json:"theme-background3"`
	OnSurface string `json:"theme-text"`
	Primary   string `json:"theme-primary"`
	Secondary string `json:"theme-secondary"`
	Accent    string `json:"theme-accent"`
}

// JSMetaData stores data about the pages.
type JSMetaData struct {
	Pages []string `json:"pages"`
	//PageToTags      map[string][]string `json:"pageTagsMap"`
	CategoryToPageIDs map[string][]int `json:"categoryToPageIDs"`
}

// A StaticSiteGenerator generates static wiki site content.
type StaticSiteGenerator struct {
	Path string

	Config *WikiConfig
	JSMeta *JSMetaData

	HTMLTemplate *template.Template
	CSSTemplate  *template.Template
	JSTemplate   *template.Template

	MDParser   goldmark.Markdown
	TitleCaser cases.Caser
}

// NewStaticSiteGenerator returns a StaticSiteGenerator.
// A _config.json file will be deserialized if present.
func NewStaticSiteGenerator(path string) *StaticSiteGenerator {
	cfg := &WikiConfig{
		Title:     "Easy Wiki",
		Surface:   "#242424",
		Surface2:  "#363636",
		Surface3:  "#484848",
		OnSurface: "#FFFFFF",
		Primary:   "#C588F9",
		Secondary: "#5E9ED6",
		Accent:    "#F6C177",
	}
	jsMeta := &JSMetaData{
		Pages:             make([]string, 0, 128),
		CategoryToPageIDs: make(map[string][]int),
	}

	cfgJSON, _ := os.ReadFile(filepath.Join(path, "_config.json"))
	_ = json.Unmarshal(cfgJSON, cfg)

	//------------------------- HTML TEMPLATE ----------------------------------
	htmlTmpl := `<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="{{.CSSPath}}">
</head>

<body>
<div class="topnav">
  <a href="{{.IndexPath}}">Home</a>
</div>
{{.Body}}
</body>
</html>`

	//------------------------- CSS TEMPLATE -----------------------------------
	cssTmpl := `:root {
  --surface: {{.Surface}};
  --surface2: {{.Surface2}};
  --surface3: {{.Surface3}};
  --onsurface: {{.OnSurface}};
  --primary: {{.Primary}};
  --secondary: {{.Secondary}};
  --accent: {{.Accent}};
}

*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: inherit;
  background-color: var(--surface);
  color: var(--onsurface);
}

*:focus {
  outline: 2px solid var(--secondary);
}

html {
  font-size: 62.5%;
}

body {
  box-sizing: border-box;
  font-size: 1.6rem;
  height: 100vh;
}

p, h1, h2, h3, a {
  opacity: 0.80;
  text-decoration: none;
  font-family: Helvetica, Verdana, sans-serif;
}

.center-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  text-align: center;
  max-width: 640px;
  margin: auto;
  padding: 0px 10px;
}

.center-content > h1 {
  padding: 114px 32px 64px 32px;
  margin-bottom: 4px;
}

#search-container {
  border-radius: 24px;
}

#search-bar {
  background-color: var(--surface3);
  border-radius: 24px;
  border: none;
  height: 48px;
  width: 100%;
  padding: 0px 24px 0px 24px;
  font-size: 2.0rem;
}

.dropdown-content {
  background-color: rgba(0,0,0,0);
  text-align: left;
  width: 100%;
  margin-top: 4px;
}

.dropdown-content > a {
  background-color: rgba(0,0,0,0);
  display: inline-block;
  width: 100%;
  padding: 8px;
  font-weight: bold;
  border-radius: 6px;
  opacity: 0.60;
  margin-bottom: 4px;
}

.dropdown-content > a:hover,
.dropdown-content > a:focus {
  opacity: 0.87;
}

.topnav {
  position: fixed;
  top: 0px;
  display: flex;
  height: 48px;
  width: 100%;
  text-align: center;
  border-bottom: 1px solid var(--surface2);
  z-index: 10;
}

.topnav > a {
  align-content: center;
  min-width: 96px;
  font-weight: bold;
  opacity: 0.60;
  margin: 4px;
  z-index: 11;
}

.topnav > a:hover,
.topnav > a:focus {
  opacity: 0.87;
}

.page-content {
  padding: 64px 16px;
  word-wrap: break-word;
  line-height: 1.5;
  max-width: 1012px;
  margin-right: auto;
  margin-left: auto;
}

.page-content > :is(h1, h2, h3, h4, h5, h6) {
  margin-top: 1.6rem;
  margin-bottom: 1.0rem;
}

.page-content > p {
  margin-top: 1.0rem;
  margin-bottom: 1.0rem;
}

.page-content a {
  color: var(--secondary);
  opacity: 1.0;
}

.page-content a:hover {
  color: var(--primary);
  text-decoration: underline;
}

blockquote {
  border-left: 6px solid var(--surface3);
  margin: 1.5em 10px;
  padding: 0.5em 10px;
}

blockquote p {
  display: inline;
}

ul, ol {
  margin: 0.5em 10px 1.0em 10px;
}

li {
  margin: 0.5em 10px 0.5em 10px ;
}

li > p {
  margin-left: 0.5em;
}

code {
  display: inline;
  font-family: "Lucida Console", Monaco, monospace;
  color: var(--accent);
  background: var(--surface2);
  border-radius: 4px;
  padding: 2px;
  margin: 0;
  overflow: visible;
  line-height: inherit;
  word-wrap: normal;
  word-break: normal;
  opacity: 1.0;
}

pre {
  display: block;
  background: var(--surface2);
  border-radius: 6px;
  padding: 16px;
  margin-top: 0;
  margin-bottom: 16px;
  overflow: auto;
  line-height: 1.45;
  word-wrap: normal;
}

img {
  max-width: 100%;
  box-sizing: content-box;
}

table {
  display: block;
  width: max-content;
  max-width: 100%;
  overflow: auto;
  margin-top: 0;
  margin-bottom: 16px;
  border-collapse: collapse;
  border-spacing: 0;
  border-color: var(--surface3);
  font-family: sans-serif;
  opacity: 0.8;
}

thead {
  display: table-header-group;
  vertical-align: middle;
}

tr {
  border-top: 1px solid var(--surface3);
  vertical-align: inherit;
  border-color: inherit;
}

tr:nth-child(even) {background-color: var(--surface2);}

th, td {
  padding: 6px 13px;
  border: 1px solid var(--surface3);
  background-color: rgba(0,0,0,0);
}

th {
  font-weight: 600;
}

td {
  display: table-cell;
  vertical-align: inherit;
}`

	//------------------------- JS TEMPLATE ------------------------------------
	jsTmpl := `class RadixNode{constructor(e,i=!1){this.edgeLabel=e,this.children=Object.create(null),this.isWord=i}static largestCommonPrefix(i,s){let t="";for(let e=0;e<Math.min(i.length,s.length);e++){if(i[e]!==s[e])return t;t+=i[e]}return t}insert(e){var i,s,t,r,h,n;this.edgeLabel!==e||this.isWord?(i=e[0])in this.children?(r=(n=this.children[i]).edgeLabel,t=(s=RadixNode.largestCommonPrefix(r,e))[0],r=r.substring(s.length),h=e.substring(s.length),""!==r&&(n.edgeLabel=r,n=this.children[t],this.children[t]=new RadixNode(s,!1),this.children[t].children[r[0]]=n,""===h)?this.children[t].isWord=!0:this.children[t].insert(h)):this.children[i]=new RadixNode(e,!0):this.isWord=!0}static searchList=[];search(e,i=""){if(i+=this.edgeLabel,e=e.substring(this.edgeLabel.length),this.isWord&&i.includes(e)&&RadixNode.searchList.push(i),""!==e)e[0]in this.children&&this.children[e[0]].search(e,i);else for(var s of Object.values(this.children))s.search("",i)}print(e=0){0!==this.edgeLabel&&console.log("-".repeat(e),this.edgeLabel,this.isWord?" (leaf)":"");for(var i of Object.values(this.children))i.print(e+1)}}class RadixTree{constructor(){this.root=new RadixNode("")}insert(e){this.root.insert(e)}search(e){var i=[];for(this.root.search(e);0<RadixNode.searchList.length;)i.push(RadixNode.searchList.pop());return i}print(){this.root.print()}}

const metaData = {{.}};
const pageTrie = new RadixTree();
const dropdown = document.getElementById("dropdown");
const searchContainer = document.getElementById("search-container");
const searchBar = document.getElementById("search-bar");

function capitalize(str) {
  return str[0].toUpperCase() + str.substring(1);
}

function nameToTitle(str) {
  return str.split('-').map(capitalize).join(' ');
}

function updateSearchResults() {
  let resultNodes = [];
  const searchText = searchBar.value.trim().toLowerCase();

  if (!searchText) {
    dropdown.replaceChildren(resultNodes);
    searchContainer.style.backgroundColor = "var(--surface)";
    return;
  }

  const pages = pageTrie.search(searchText);
  for (const page of pages) {
    let result = document.createElement('a');
    result.innerText = nameToTitle(page);
    result.className = "search-result";
    result.href = "static/pages/" + page + ".html";
    resultNodes.push(result);
  }
  if (resultNodes.length > 0) {
    dropdown.replaceChildren(...resultNodes);
    searchContainer.style.backgroundColor = "var(--surface3)";
  }
}

console.log("Building page index..");
for (const page of metaData.pages) {
  pageTrie.insert(page);
}
console.log("Done.");`

	return &StaticSiteGenerator{
		Path:         path,
		Config:       cfg,
		JSMeta:       jsMeta,
		HTMLTemplate: template.Must(template.New("layout").Parse(htmlTmpl)),
		CSSTemplate:  template.Must(template.New("styles").Parse(cssTmpl)),
		JSTemplate:   template.Must(template.New("script").Parse(jsTmpl)),
		MDParser: goldmark.New(
			goldmark.WithExtensions(extension.GFM, meta.Meta, extension.Typographer),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
		TitleCaser: cases.Title(language.English),
	}
}

// Initialize generates the wiki directory structure and static content.
func (ssg StaticSiteGenerator) Initialize() {
	_ = CreateDir(ssg.Path)
	ssg.MakeConfig()
	ssg.MakeWiki()
}

func (ssg StaticSiteGenerator) RenderAll() {
	ssg.RenderIndex()
	ssg.RenderPages()
	ssg.RenderCSS()
	ssg.RenderJS()
}

// MakeConfig attempts to create the _config.json configuration file.
func (ssg StaticSiteGenerator) MakeConfig() {
	cfgJSON, err := json.MarshalIndent(ssg.Config, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	cfgFile, err := os.Create(filepath.Join(ssg.Path, "_config.json"))
	if err != nil {
		log.Println(err)
		return
	}
	defer cfgFile.Close()

	_, err = cfgFile.Write(cfgJSON)
	if err != nil {
		log.Println(err)
	}
}

// MakeWiki generates the necessary files and directory structure for the wiki.
func (ssg StaticSiteGenerator) MakeWiki() {
	_ = CreateDir(filepath.Join(ssg.Path, "_pages"))
	_ = CreateDir(filepath.Join(ssg.Path, "static"))
	_ = CreateDir(filepath.Join(ssg.Path, "static/pages"))

	ssg.RenderIndex()
	ssg.RenderCSS()
	ssg.RenderJS()
}

func (ssg StaticSiteGenerator) RenderIndex() {

	indexBody := "<div class=\"center-content\">\n" +
		"  <h1>" + ssg.Config.Title + "</h1>" + `
  <div id="search-container">
    <input type="text" id="search-bar" placeholder="Search..." onkeyup="updateSearchResults();" autocomplete="off" autofocus="true">
    <div id="dropdown" class="dropdown-content">
    </div>
  </div>
</div>
<script src="bundle.js"></script>`

	indexFile, err := os.Create(filepath.Join(ssg.Path, "index.html"))
	if err != nil {
		log.Println(err)
		return
	}
	defer indexFile.Close()

	err = ssg.HTMLTemplate.Execute(indexFile, struct {
		Title     string
		CSSPath   string
		IndexPath string
		Body      string
	}{
		Title:     ssg.Config.Title,
		CSSPath:   "styles.css",
		IndexPath: "index.html",
		Body:      indexBody,
	})
	if err != nil {
		log.Println(err)
	}
}

func (ssg StaticSiteGenerator) RenderPages() {
	var htmlBuf bytes.Buffer

	mdPagesDir := filepath.Join(ssg.Path, "_pages")
	htmlPagesDir := filepath.Join(ssg.Path, "static/pages")

	mdPages, err := os.ReadDir(mdPagesDir)
	if err != nil {
		log.Println(err)
		return
	}

	for pageID, mdPage := range mdPages {
		mdNameExt := strings.Split(mdPage.Name(), ".")
		mdName := strings.ToLower(mdNameExt[0])
		mdExt := strings.Join(mdNameExt[1:], "")

		if mdExt != "md" {
			continue
		}

		ssg.JSMeta.Pages = append(ssg.JSMeta.Pages, mdName)

		mdData, err := os.ReadFile(filepath.Join(mdPagesDir, mdPage.Name()))
		if err != nil {
			log.Println(err)
			continue
		}

		htmlFile, err := os.Create(filepath.Join(htmlPagesDir, mdName+".html"))
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = htmlBuf.WriteString("<div class=\"page-content\">\n")
		if err != nil {
			log.Println(err)
		}

		ctx := parser.NewContext()
		err = ssg.MDParser.Convert(mdData, &htmlBuf, parser.WithContext(ctx))
		if err != nil {
			log.Println(err)
		}

		_, err = htmlBuf.WriteString("</div>")
		if err != nil {
			log.Println(err)
		}

		metaData := meta.Get(ctx)
		title, ok := metaData["Title"].(string)
		category, ok := metaData["Category"].(string)

		if title == "" || !ok {
			title = ssg.TitleCaser.String(mdName)
		}

		if category != "" && ok {
			if _, ok := ssg.JSMeta.CategoryToPageIDs[category]; !ok {
				ssg.JSMeta.CategoryToPageIDs[category] = make([]int, 0)
			}
			ssg.JSMeta.CategoryToPageIDs[category] = append(ssg.JSMeta.CategoryToPageIDs[category], pageID)
		}

		err = ssg.HTMLTemplate.Execute(htmlFile, struct {
			Title     string
			CSSPath   string
			IndexPath string
			Body      string
		}{
			Title:     title,
			CSSPath:   filepath.Join("..", "..", "styles.css"),
			IndexPath: filepath.Join("..", "..", "index.html"),
			Body:      strings.TrimSpace(htmlBuf.String()),
		})
		if err != nil {
			log.Println(err)
		}

		htmlBuf.Reset()
		htmlFile.Close()
	}
}

func (ssg StaticSiteGenerator) RenderJS() {
	jsFile, err := os.Create(filepath.Join(ssg.Path, "bundle.js"))
	if err != nil {
		log.Println(err)
		return
	}
	defer jsFile.Close()

	jsMeta, err := json.Marshal(ssg.JSMeta)
	if err != nil {
		log.Println(err)
		return
	}

	err = ssg.JSTemplate.Execute(jsFile, string(jsMeta))
	if err != nil {
		log.Println(err)
	}
}

func (ssg StaticSiteGenerator) RenderCSS() {
	cssFile, err := os.Create(filepath.Join(ssg.Path, "styles.css"))
	if err != nil {
		log.Println(err)
		return
	}
	defer cssFile.Close()

	err = ssg.CSSTemplate.Execute(cssFile, ssg.Config)
	if err != nil {
		log.Println(err)
	}
}

func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
