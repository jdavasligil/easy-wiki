import MarkdownIt from "markdown-it";
import fs from 'node:fs';

const pageDir = 'pages';
const htmlDir = 'wiki/generated';
const indexFile = 'wiki/index.js';
const pageTemplate = fs.readFileSync(`page-template.html`, 'utf8');

const helpText = `
Usage: server.js [<optional-argument>]
  options:
    -r    Re-render all pages on startup
`;

const md = MarkdownIt({
  // Enable HTML tags in source
  html: true,

  // Use '/' to close single tags (<br />).
  // This is only for full CommonMark compatibility.
  xhtmlOut: false,

  // Convert '\n' in paragraphs into <br>
  breaks: false,

  // CSS language prefix for fenced blocks. Can be
  // useful for external highlighters.
  langPrefix: 'language-',

  // Autoconvert URL-like text to links
  linkify: true,

  // Enable some language-neutral replacement + quotes beautification
  // For the full list of replacements, see https://github.com/markdown-it/markdown-it/blob/master/lib/rules_core/replacements.mjs
  typographer:  true,

  // Double + single quotes replacement pairs, when typographer enabled,
  // and smartquotes on. Could be either a String or an Array.
  //
  // For example, you can use '«»„“' for Russian, '„“‚‘' for German,
  // and ['«\xA0', '\xA0»', '‹\xA0', '\xA0›'] for French (including nbsp).
  quotes: '“”‘’',

  // Highlighter function. Should return escaped HTML,
  // or '' if the source string is not changed and should be escaped externally.
  // If result starts with <pre... internal wrapper is skipped.
  highlight: function (/*str, lang*/) { return ''; }
});

function readPages() {
  try {
    return fs.readdirSync(pageDir);
  } catch (err) {
    console.error(err);
    return [];
  }
}

function renderPage(pagename) {
  try {
    const mdFile = fs.readFileSync(`${pageDir}/${pagename}.md`, 'utf8');
    const titleIdx = pageTemplate.indexOf('</title>');
    const bodyIdx = pageTemplate.indexOf('</body>');
    const html = pageTemplate.substring(0, titleIdx) + pagename + pageTemplate.substring(titleIdx, bodyIdx) +  md.render(mdFile) + pageTemplate.substring(bodyIdx);
    fs.writeFileSync(`${htmlDir}/${pagename}.html`, html);
  } catch (err) {
    console.error(err);
  }
}

function renderAll() {
  const pages = readPages();
  for (const page of pages) {
    const [ pagename, ext ] = page.split(".");
    if (ext === 'md') {
      renderPage(pagename);
    }
  }
}

function deletePage(pagename) {
  try {
    fs.unlinkSync(`${htmlDir}/${pagename}.html`);
  } catch (err) {
    console.error(err);
  }
}


// TODO: Refine to replace generated area. Include a map for category => pageIDs.
function updateIndex() {
  const pages = readPages()
    .filter(page => page.endsWith('.md'))
    .map(page => '"' + page.split('.')[0] + '"')
    .join(',');
  console.log(pages);

  const indexJS = fs.readFileSync(indexFile, 'utf8');
  const pageIdx = indexJS.indexOf('[');
  const pageEndIdx = indexJS.indexOf(']');
  fs.writeFileSync(indexFile, `${indexJS.substring(0, pageIdx + 1)}${pages}${indexJS.substring(pageEndIdx)}`.trim());
}


// Handle command line arguments.
if (process.argv[2]) {
  const opt = process.argv[2].replace(/-*/, '');
  switch (opt) {
  case 'r':
    renderAll();
    updateIndex();
    break;
  default:
    console.log(helpText);
  }
}

// Begin file watching for markdown updates.
fs.watch(pageDir, (eventType, filename) => {
  if (!filename)
    return;

  const [ pagename, ext ] = filename.split(".");

  if (!pagename || ext !== "md")
    return;

  if (eventType === "change") {
    console.log(`filename: ${filename}`);
    renderPage(pagename);
    console.log("HTML Updated");
    updateIndex();
  } else if (!fs.existsSync(`${pageDir}/${filename}`)) {
    console.log(`filename: ${filename}`);
    deletePage(pagename);
    console.log("Page Deleted");
    updateIndex();
  }
});
