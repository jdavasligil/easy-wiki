package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// A WikiConfig is a (de)serializable config which holds extensible settings.
type WikiConfig struct {
	Title    string `json:"title"`
	Darkmode bool   `json:"darkmode"`
}

// MetaData stores data about the pages.
type MetaData struct {
	Pages []string `json:"pages"`
	//PageToTags      map[string][]string `json:"pageTagsMap"`
	CategoryToPages map[string][]string `json:"categoryToPages"`
}

// A StaticSiteGenerator generates static wiki site content.
type StaticSiteGenerator struct {
	Config       *WikiConfig
	Meta         *MetaData
	HTMLTemplate *template.Template
	JSTemplate   *template.Template
	Path         string
}

// NewStaticSiteGenerator returns a StaticSiteGenerator.
// A _config.json file will be deserialized if present.
func NewStaticSiteGenerator(path string) *StaticSiteGenerator {
	cfg := &WikiConfig{}
	meta := &MetaData{Pages: []string{"al", "an", "andrew"}}
	cfgJSON, _ := os.ReadFile(filepath.Join(path, "_config.json"))
	_ = json.Unmarshal(cfgJSON, cfg)

	htmlTmpl := `<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="{{.CSSPath}}">
</head>

<body>
{{.Body}}
</body>
</html>`

	jsTmpl := `// GENERATED (DO NOT EDIT)
const metaData = {{.}};
// END GENERATED

class RadixNode{constructor(e,i=!1){this.edgeLabel=e,this.children=Object.create(null),this.isWord=i}static largestCommonPrefix(i,s){let t="";for(let e=0;e<Math.min(i.length,s.length);e++){if(i[e]!==s[e])return t;t+=i[e]}return t}insert(e){var i,s,t,r,h,n;this.edgeLabel!==e||this.isWord?(i=e[0])in this.children?(r=(n=this.children[i]).edgeLabel,t=(s=RadixNode.largestCommonPrefix(r,e))[0],r=r.substring(s.length),h=e.substring(s.length),""!==r&&(n.edgeLabel=r,n=this.children[t],this.children[t]=new RadixNode(s,!1),this.children[t].children[r[0]]=n,""===h)?this.children[t].isWord=!0:this.children[t].insert(h)):this.children[i]=new RadixNode(e,!0):this.isWord=!0}static searchList=[];search(e,i=""){if(i+=this.edgeLabel,e=e.substring(this.edgeLabel.length),this.isWord&&i.includes(e)&&RadixNode.searchList.push(i),""!==e)e[0]in this.children&&this.children[e[0]].search(e,i);else for(var s of Object.values(this.children))s.search("",i)}print(e=0){0!==this.edgeLabel&&console.log("-".repeat(e),this.edgeLabel,this.isWord?" (leaf)":"");for(var i of Object.values(this.children))i.print(e+1)}}class RadixTree{constructor(){this.root=new RadixNode("")}insert(e){this.root.insert(e)}search(e){var i=[];for(this.root.search(e);0<RadixNode.searchList.length;)i.push(RadixNode.searchList.pop());return i}print(){this.root.print()}}

const pageTrie = new RadixTree();
const dropdown = document.getElementById("dropdown");
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
    return;
  }

  const pages = pageTrie.search(searchText);
  for (const page of metaData.pages) {
    let result = document.createElement('a');
    result.innerText = nameToTitle(page);
    result.className = "search-result";
    result.href = "generated/" + page + ".html";
    resultNodes.push(result);
  }
  dropdown.replaceChildren(...resultNodes);
}

console.log("Building page index..");
for (const page of metaData.pages) {
  pageTrie.insert(page);
}
console.log("Done.");`

	return &StaticSiteGenerator{
		Config:       cfg,
		Meta:         meta,
		HTMLTemplate: template.Must(template.New("layout").Parse(htmlTmpl)),
		JSTemplate:   template.Must(template.New("script").Parse(jsTmpl)),
		Path:         path,
	}
}

// Initialize generates the wiki directory structure and static content.
func (ssg StaticSiteGenerator) Initialize() {
	_ = CreateDir(ssg.Path)
	ssg.MakeConfig()
	ssg.MakeWiki()
}

// MakeConfig attempts to create the _config.json configuration file.
func (ssg StaticSiteGenerator) MakeConfig() {
	cfgJSON, err := json.Marshal(ssg.Config)
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

func (ssg StaticSiteGenerator) MakeWiki() {
	_ = CreateDir(filepath.Join(ssg.Path, "_pages"))
	_ = CreateDir(filepath.Join(ssg.Path, "static"))
	_ = CreateDir(filepath.Join(ssg.Path, "static/pages"))

	ssg.RenderIndex()
	ssg.RenderJS()
}

func (ssg StaticSiteGenerator) RenderIndex() {
	indexBody := "<h1>" + ssg.Config.Title + "</h1>" + `
  <div class="search-container">
    <input type="search" id="search-bar" placeholder="Search..." onkeyup="updateSearchResults();" autocomplete="off" autofocus="true">
    <div id="dropdown" class="dropdown-content">
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
		Title   string
		CSSPath string
		Body    string
	}{
		Title:   ssg.Config.Title,
		CSSPath: "styles.css",
		Body:    indexBody,
	})
	if err != nil {
		log.Println(err)
	}
}

func (ssg StaticSiteGenerator) RenderJS() {
	jsFile, err := os.Create(filepath.Join(ssg.Path, "bundle.js"))
	if err != nil {
		log.Println(err)
		return
	}
	defer jsFile.Close()

	metaData, err := json.Marshal(ssg.Meta)
	if err != nil {
		log.Println(err)
		return
	}

	err = ssg.JSTemplate.Execute(jsFile, string(metaData))
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
