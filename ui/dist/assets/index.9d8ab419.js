import{r as e,o as t,c as r,a as o,b as n,d as i,i as s,e as l}from"./vendor.7e9e41ae.js";!function(){const e=document.createElement("link").relList;if(!(e&&e.supports&&e.supports("modulepreload"))){for(const e of document.querySelectorAll('link[rel="modulepreload"]'))t(e);new MutationObserver((e=>{for(const r of e)if("childList"===r.type)for(const e of r.addedNodes)"LINK"===e.tagName&&"modulepreload"===e.rel&&t(e)})).observe(document,{childList:!0,subtree:!0})}function t(e){if(e.ep)return;e.ep=!0;const t=function(e){const t={};return e.integrity&&(t.integrity=e.integrity),e.referrerpolicy&&(t.referrerPolicy=e.referrerpolicy),"use-credentials"===e.crossorigin?t.credentials="include":"anonymous"===e.crossorigin?t.credentials="omit":t.credentials="same-origin",t}(e);fetch(e.href,t)}}();const c={};c.render=function(o,n){const i=e("router-view");return t(),r(i)};var a=Object.freeze({__proto__:null,[Symbol.toStringTag]:"Module",default:c});const d={},u=function(e,t){return t&&0!==t.length?Promise.all(t.map((e=>{if((e=`/${e}`)in d)return;d[e]=!0;const t=e.endsWith(".css"),r=t?'[rel="stylesheet"]':"";if(document.querySelector(`link[href="${e}"]${r}`))return;const o=document.createElement("link");return o.rel=t?"stylesheet":"modulepreload",t||(o.as="script",o.crossOrigin=""),o.href=e,document.head.appendChild(o),t?new Promise(((e,t)=>{o.addEventListener("load",e),o.addEventListener("error",t)})):void 0}))).then((()=>e())):e()},m=[{path:"/",redirect:"/ui"},{path:"/ui",component:()=>u((()=>Promise.resolve().then((function(){return a}))),void 0),redirect:"/ui/index",children:[{path:"index",component:()=>u((()=>import("./index.4b987154.js")),["assets/index.4b987154.js","assets/index.7290fb28.css","assets/index.2cf0d985.js","assets/vendor.7e9e41ae.js"]),meta:{title:"Bridge Test"}},{path:"lora",component:()=>u((()=>import("./lora.39b9661e.js")),["assets/lora.39b9661e.js","assets/lora.9dfdb658.css","assets/index.2cf0d985.js","assets/vendor.7e9e41ae.js"]),meta:{title:"Lora-ttn Test"}}]}],p=o({history:n("#"),routes:m});p.beforeEach(((e,t,r)=>{e.meta.title&&(document.title=e.meta.title),r()}));const f=i(c);l.pop({title:"Example",css:{width:"250px",height:"180px"},showMask:!0}),f.use(s),f.use(p),f.mount("#app");
