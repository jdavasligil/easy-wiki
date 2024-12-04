class RadixNode{constructor(e,i=!1){this.edgeLabel=e,this.children=Object.create(null),this.isWord=i}static largestCommonPrefix(i,s){let t="";for(let e=0;e<Math.min(i.length,s.length);e++){if(i[e]!==s[e])return t;t+=i[e]}return t}insert(e){var i,s,t,r,h,n;this.edgeLabel!==e||this.isWord?(i=e[0])in this.children?(r=(n=this.children[i]).edgeLabel,t=(s=RadixNode.largestCommonPrefix(r,e))[0],r=r.substring(s.length),h=e.substring(s.length),""!==r&&(n.edgeLabel=r,n=this.children[t],this.children[t]=new RadixNode(s,!1),this.children[t].children[r[0]]=n,""===h)?this.children[t].isWord=!0:this.children[t].insert(h)):this.children[i]=new RadixNode(e,!0):this.isWord=!0}static searchList=[];search(e,i=""){if(i+=this.edgeLabel,e=e.substring(this.edgeLabel.length),this.isWord&&i.includes(e)&&RadixNode.searchList.push(i),""!==e)e[0]in this.children&&this.children[e[0]].search(e,i);else for(var s of Object.values(this.children))s.search("",i)}print(e=0){0!==this.edgeLabel&&console.log("-".repeat(e),this.edgeLabel,this.isWord?" (leaf)":"");for(var i of Object.values(this.children))i.print(e+1)}}class RadixTree{constructor(){this.root=new RadixNode("")}insert(e){this.root.insert(e)}search(e){var i=[];for(this.root.search(e);0<RadixNode.searchList.length;)i.push(RadixNode.searchList.pop());return i}print(){this.root.print()}}

const pageTrie = new RadixTree();
const pages = document.getElementById("pages").innerText;
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
  for (const page of pages) {
    let result = document.createElement('a');
    result.innerText = nameToTitle(page);
    result.className = "search-result";
    result.href = `generated/${page}.html`;
    resultNodes.push(result);
  }
  dropdown.replaceChildren(...resultNodes);
}

console.log("Building page index..");
for (const page of pages.split(' ')) {
  pageTrie.insert(page);
}
console.log("Done.");
