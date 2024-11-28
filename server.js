import MarkdownIt from "markdown-it";
import fs from 'node:fs'

const md = MarkdownIt({
  // Enable HTML tags in source
  html: true,

  // Use '/' to close single tags (<br />).
  // This is only for full CommonMark compatibility.
  xhtmlOut: false,

  // Convert '\n' in paragraphs into <br>
  breaks: true,

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
    return fs.readdirSync('pages');
  } catch (err) {
    console.error(err);
    return [];
  }
}

// Loop through .md files and render to html
//(() => {
//  const pages = readPages();
//  console.log(pages);
//  try {
//    const data = fs.readFileSync('pages/test.md', 'utf8');
//    console.log(md.render(data));
//  } catch (err) {
//    console.error(err);
//  }
//})();

fs.watch('pages', (eventType, filename) => {
  if (!filename)
    return;

  const [ base, ext ] = filename.split(".");

  if (!base || ext !== "md")
    return;

  if (eventType === "change") {
    console.log(`filename: ${filename}`);
    console.log("Update HTML");
  } else if (!fs.existsSync(`pages/${filename}`)) {
    console.log(`filename: ${filename}`);
    console.log("Page Deleted");
  }
});
