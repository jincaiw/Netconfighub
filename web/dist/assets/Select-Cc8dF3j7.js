import{m as e}from"./api-CZuskCjY.js";import{$t as t,An as n,Ct as r,En as i,Et as a,Lt as o,Pt as s,Q as c,Tt as l,Vt as u,X as d,Yt as f,Z as p,Zt as m,_n as h,an as g,at as _,bt as v,c as y,dt as b,en as x,g as S,gt as C,h as w,hn as T,ht as E,in as D,jt as O,ln as k,m as A,mn as j,mt as ee,nn as M,on as N,ot as P,pt as F,rt as I,s as L,tn as R,ut as z,wt as B,xt as V,y as H}from"./Button-v7qKim5p.js";import{C as te,S as U,_ as W,b as ne,c as G,d as K,f as q,g as J,h as re,m as ie,o as ae,p as oe,r as se,s as ce,u as le,v as Y,y as X}from"./Tooltip-qMw84P-o.js";import{D as Z,E as Q,_ as ue,f as de,h as fe,k as pe,m as me,u as he,v as ge,w as _e}from"./token-dQuugHoF.js";import{o as ve,r as ye}from"./Input-C6Yqttcz.js";import{S as be,f as $,s as xe,v as Se}from"./index-CESIsNoN.js";function Ce(e){return e&-e}var we=class{constructor(e,t){this.l=e,this.min=t;let n=Array(e+1);for(let t=0;t<e+1;++t)n[t]=0;this.ft=n}add(e,t){if(t===0)return;let{l:n,ft:r}=this;for(e+=1;e<=n;)r[e]+=t,e+=Ce(e)}get(e){return this.sum(e+1)-this.sum(e)}sum(e){if(e===void 0&&(e=this.l),e<=0)return 0;let{ft:t,min:n,l:r}=this;if(e>r)throw Error("[FinweckTree.sum]: `i` is larger than length.");let i=e*n;for(;e>0;)i+=t[e],e-=Ce(e);return i}getBound(e){let t=0,n=this.l;for(;n>t;){let r=Math.floor((t+n)/2),i=this.sum(r);if(i>e){n=r;continue}else if(i<e){if(t===r)return this.sum(t+1)<=e?t+1:r;t=r}else return r}return t}},Te;function Ee(){return typeof document>`u`?!1:(Te===void 0&&(Te=`matchMedia`in window?window.matchMedia(`(pointer:coarse)`).matches:!1),Te)}var De;function Oe(){return typeof document>`u`?1:(De===void 0&&(De=`chrome`in window?window.devicePixelRatio:1),De)}var ke=`VVirtualListXScroll`;function Ae({columnsRef:e,renderColRef:t,renderItemWithColsRef:n}){let r=i(0),a=i(0),o=u(()=>{let t=e.value;if(t.length===0)return null;let n=new we(t.length,0);return t.forEach((e,t)=>{n.add(t,e.width)}),n});return k(ke,{startIndexRef:E(()=>{let e=o.value;return e===null?0:Math.max(e.getBound(a.value)-1,0)}),endIndexRef:E(()=>{let t=o.value;return t===null?0:Math.min(t.getBound(a.value+r.value)+1,e.value.length-1)}),columnsRef:e,renderColRef:t,renderItemWithColsRef:n,getLeft:e=>{let t=o.value;return t===null?0:t.sum(e)}}),{listWidthRef:r,scrollLeftRef:a}}var je=f({name:`VirtualListRow`,props:{index:{type:Number,required:!0},item:{type:Object,required:!0}},setup(){let{startIndexRef:e,endIndexRef:n,columnsRef:r,getLeft:i,renderColRef:a,renderItemWithColsRef:o}=t(ke);return{startIndex:e,endIndex:n,columns:r,renderCol:a,renderItemWithCols:o,getLeft:i}},render(){let{startIndex:e,endIndex:t,columns:n,renderCol:r,renderItemWithCols:i,getLeft:a,item:o}=this;if(i!=null)return i({itemIndex:this.index,startColIndex:e,endColIndex:t,allColumns:n,item:o,getLeft:a});if(r!=null){let i=[];for(let s=e;s<=t;++s){let e=n[s];i.push(r({column:e,left:a(s),item:o}))}return i}return null}}),Me=oe(`.v-vl`,{maxHeight:`inherit`,height:`100%`,overflow:`auto`,minWidth:`1px`},[oe(`&:not(.v-vl--show-scrollbar)`,{scrollbarWidth:`none`},[oe(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,{width:0,height:0,display:`none`})])]),Ne=f({name:`VirtualList`,inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},columns:{type:Array,default:()=>[]},renderCol:Function,renderItemWithCols:Function,items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:`div`},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:`key`},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(e){let t=b();Me.mount({id:`vueuc/virtual-list`,head:!0,anchorMetaName:ie,ssr:t}),N(()=>{let{defaultScrollIndex:t,defaultScrollKey:n}=e;t==null?n!=null&&x({key:n}):x({index:t})});let r=!1,a=!1;M(()=>{if(r=!1,!a){a=!0;return}x({top:_.value,left:c.value})}),g(()=>{r=!0,a||=!0});let o=E(()=>{if(e.renderCol==null&&e.renderItemWithCols==null||e.columns.length===0)return;let t=0;return e.columns.forEach(e=>{t+=e.width}),t}),s=u(()=>{let t=new Map,{keyField:n}=e;return e.items.forEach((e,r)=>{t.set(e[n],r)}),t}),{scrollLeftRef:c,listWidthRef:l}=Ae({columnsRef:n(e,`columns`),renderColRef:n(e,`renderCol`),renderItemWithColsRef:n(e,`renderItemWithCols`)}),d=i(null),f=i(void 0),p=new Map,m=u(()=>{let{items:t,itemSize:n,keyField:r}=e,i=new we(t.length,n);return t.forEach((e,t)=>{let n=e[r],a=p.get(n);a!==void 0&&i.add(t,a)}),i}),h=i(0),_=i(0),v=E(()=>Math.max(m.value.getBound(_.value-_e(e.paddingTop))-1,0)),y=u(()=>{let{value:t}=f;if(t===void 0)return[];let{items:n,itemSize:r}=e,i=v.value,a=Math.min(i+Math.ceil(t/r+1),n.length-1),o=[];for(let e=i;e<=a;++e)o.push(n[e]);return o}),x=(e,t)=>{if(typeof e==`number`){T(e,t,`auto`);return}let{left:n,top:r,index:i,key:a,position:o,behavior:c,debounce:l=!0}=e;if(n!==void 0||r!==void 0)T(n,r,c);else if(i!==void 0)w(i,c,l);else if(a!==void 0){let e=s.value.get(a);e!==void 0&&w(e,c,l)}else o===`bottom`?T(0,2**53-1,c):o===`top`&&T(0,0,c)},S,C=null;function w(t,n,r){let{value:i}=m,a=i.sum(t)+_e(e.paddingTop);if(!r)d.value.scrollTo({left:0,top:a,behavior:n});else{S=t,C!==null&&window.clearTimeout(C),C=window.setTimeout(()=>{S=void 0,C=null},16);let{scrollTop:e,offsetHeight:r}=d.value;if(a>e){let o=i.get(t);a+o<=e+r||d.value.scrollTo({left:0,top:a+o-r,behavior:n})}else d.value.scrollTo({left:0,top:a,behavior:n})}}function T(e,t,n){d.value.scrollTo({left:e,top:t,behavior:n})}function D(t,n){if(r||e.ignoreItemResize||F(n.target))return;let{value:i}=m,a=s.value.get(t),o=i.get(a),c=n.borderBoxSize?.[0]?.blockSize??n.contentRect.height;if(c===o)return;c-e.itemSize===0?p.delete(t):p.set(t,c-e.itemSize);let l=c-o;if(l===0)return;i.add(a,l);let u=d.value;if(u!=null){if(S===void 0){let e=i.sum(a);u.scrollTop>e&&u.scrollBy(0,l)}else (a<S||a===S&&c+i.sum(a)>u.scrollTop+u.offsetHeight)&&u.scrollBy(0,l);P()}h.value++}let O=!Ee(),k=!1;function A(t){var n;(n=e.onScroll)==null||n.call(e,t),(!O||!k)&&P()}function j(t){var n;if((n=e.onWheel)==null||n.call(e,t),O){let e=d.value;if(e!=null){if(t.deltaX===0&&(e.scrollTop===0&&t.deltaY<=0||e.scrollTop+e.offsetHeight>=e.scrollHeight&&t.deltaY>=0))return;t.preventDefault(),e.scrollTop+=t.deltaY/Oe(),e.scrollLeft+=t.deltaX/Oe(),P(),k=!0,te(()=>{k=!1})}}}function ee(t){if(r||F(t.target))return;if(e.renderCol==null&&e.renderItemWithCols==null){if(t.contentRect.height===f.value)return}else if(t.contentRect.height===f.value&&t.contentRect.width===l.value)return;f.value=t.contentRect.height,l.value=t.contentRect.width;let{onResize:n}=e;n!==void 0&&n(t)}function P(){let{value:e}=d;e!=null&&(_.value=e.scrollTop,c.value=e.scrollLeft)}function F(e){let t=e;for(;t!==null;){if(t.style.display===`none`)return!0;t=t.parentElement}return!1}return{listHeight:f,listStyle:{overflow:`auto`},keyToIndex:s,itemsStyle:u(()=>{let{itemResizable:t}=e,n=Z(m.value.sum());return h.value,[e.itemsStyle,{boxSizing:`content-box`,width:Z(o.value),height:t?``:n,minHeight:t?n:``,paddingTop:Z(e.paddingTop),paddingBottom:Z(e.paddingBottom)}]}),visibleItemsStyle:u(()=>(h.value,{transform:`translateY(${Z(m.value.sum(v.value))})`})),viewportItems:y,listElRef:d,itemsElRef:i(null),scrollTo:x,handleListResize:ee,handleListScroll:A,handleListWheel:j,handleItemResize:D}},render(){let{itemResizable:e,keyField:t,keyToIndex:n,visibleItemsTag:r}=this;return m(ue,{onResize:this.handleListResize},{default:()=>{var i;return m(`div`,x(this.$attrs,{class:[`v-vl`,this.showScrollbar&&`v-vl--show-scrollbar`],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:`listElRef`}),[this.items.length===0?(i=this.$slots).empty?.call(i):m(`div`,{ref:`itemsElRef`,class:`v-vl-items`,style:this.itemsStyle},[m(r,Object.assign({class:`v-vl-visible-items`,style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>{let{renderCol:r,renderItemWithCols:i}=this;return this.viewportItems.map(a=>{let o=a[t],s=n.get(o),c=r==null?void 0:m(je,{index:s,item:a}),l=i==null?void 0:m(je,{index:s,item:a}),u=this.$slots.default({item:a,renderedCols:c,renderedItemWithCols:l,index:s})[0];return e?m(ue,{key:o,onResize:e=>this.handleItemResize(o,e)},{default:()=>u}):(u.key=o,u)})}})])])}})}});function Pe(e,t){t&&(N(()=>{let{value:n}=e;n&&ge.registerHandler(n,t)}),j(e,(e,t)=>{t&&ge.unregisterHandler(t)},{deep:!1}),D(()=>{let{value:t}=e;t&&ge.unregisterHandler(t)}))}function Fe(e,t){if(!e)return;let n=document.createElement(`a`);n.href=e,t!==void 0&&(n.download=t),document.body.appendChild(n),n.click(),document.body.removeChild(n)}function Ie(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var Le=f({name:`Checkmark`,render(){return m(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},m(`g`,{fill:`none`},m(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Re=f({name:`Empty`,render(){return m(`svg`,{viewBox:`0 0 28 28`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},m(`path`,{d:`M26 7.5C26 11.0899 23.0899 14 19.5 14C15.9101 14 13 11.0899 13 7.5C13 3.91015 15.9101 1 19.5 1C23.0899 1 26 3.91015 26 7.5ZM16.8536 4.14645C16.6583 3.95118 16.3417 3.95118 16.1464 4.14645C15.9512 4.34171 15.9512 4.65829 16.1464 4.85355L18.7929 7.5L16.1464 10.1464C15.9512 10.3417 15.9512 10.6583 16.1464 10.8536C16.3417 11.0488 16.6583 11.0488 16.8536 10.8536L19.5 8.20711L22.1464 10.8536C22.3417 11.0488 22.6583 11.0488 22.8536 10.8536C23.0488 10.6583 23.0488 10.3417 22.8536 10.1464L20.2071 7.5L22.8536 4.85355C23.0488 4.65829 23.0488 4.34171 22.8536 4.14645C22.6583 3.95118 22.3417 3.95118 22.1464 4.14645L19.5 6.79289L16.8536 4.14645Z`,fill:`currentColor`}),m(`path`,{d:`M25 22.75V12.5991C24.5572 13.0765 24.053 13.4961 23.5 13.8454V16H17.5L17.3982 16.0068C17.0322 16.0565 16.75 16.3703 16.75 16.75C16.75 18.2688 15.5188 19.5 14 19.5C12.4812 19.5 11.25 18.2688 11.25 16.75L11.2432 16.6482C11.1935 16.2822 10.8797 16 10.5 16H4.5V7.25C4.5 6.2835 5.2835 5.5 6.25 5.5H12.2696C12.4146 4.97463 12.6153 4.47237 12.865 4H6.25C4.45507 4 3 5.45507 3 7.25V22.75C3 24.5449 4.45507 26 6.25 26H21.75C23.5449 26 25 24.5449 25 22.75ZM4.5 22.75V17.5H9.81597L9.85751 17.7041C10.2905 19.5919 11.9808 21 14 21L14.215 20.9947C16.2095 20.8953 17.842 19.4209 18.184 17.5H23.5V22.75C23.5 23.7165 22.7165 24.5 21.75 24.5H6.25C5.2835 24.5 4.5 23.7165 4.5 22.75Z`,fill:`currentColor`}))}}),ze=f({props:{onFocus:Function,onBlur:Function},setup(e){return()=>m(`div`,{style:`width: 0; height: 0`,tabindex:0,onFocus:e.onFocus,onBlur:e.onBlur})}}),Be={iconSizeTiny:`28px`,iconSizeSmall:`34px`,iconSizeMedium:`40px`,iconSizeLarge:`46px`,iconSizeHuge:`52px`};function Ve(e){let{textColorDisabled:t,iconColor:n,textColor2:r,fontSizeTiny:i,fontSizeSmall:a,fontSizeMedium:o,fontSizeLarge:s,fontSizeHuge:c}=e;return Object.assign(Object.assign({},Be),{fontSizeTiny:i,fontSizeSmall:a,fontSizeMedium:o,fontSizeLarge:s,fontSizeHuge:c,textColor:t,iconColor:n,extraTextColor:r})}var He={name:`Empty`,common:L,self:Ve},Ue=V(`empty`,`
 display: flex;
 flex-direction: column;
 align-items: center;
 font-size: var(--n-font-size);
`,[r(`icon`,`
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 font-size: var(--n-icon-size);
 line-height: var(--n-icon-size);
 color: var(--n-icon-color);
 transition:
 color .3s var(--n-bezier);
 `,[v(`+`,[r(`description`,`
 margin-top: 8px;
 `)])]),r(`description`,`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 `),r(`extra`,`
 text-align: center;
 transition: color .3s var(--n-bezier);
 margin-top: 12px;
 color: var(--n-extra-text-color);
 `)]),We=f({name:`Empty`,props:Object.assign(Object.assign({},S.props),{description:String,showDescription:{type:Boolean,default:!0},showIcon:{type:Boolean,default:!0},size:{type:String,default:`medium`},renderIcon:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:r}=c(e),i=S(`Empty`,`-empty`,Ue,He,e,t),{localeRef:o}=ve(`Empty`),s=u(()=>e.description??r?.value?.Empty?.description),l=u(()=>r?.value?.Empty?.renderIcon||(()=>m(Re,null))),d=u(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{[a(`iconSize`,t)]:r,[a(`fontSize`,t)]:o,textColor:s,iconColor:c,extraTextColor:l}}=i.value;return{"--n-icon-size":r,"--n-font-size":o,"--n-bezier":n,"--n-text-color":s,"--n-icon-color":c,"--n-extra-text-color":l}}),f=n?p(`empty`,u(()=>{let t=``,{size:n}=e;return t+=n[0],t}),d,e):void 0;return{mergedClsPrefix:t,mergedRenderIcon:l,localizedDescription:u(()=>s.value||o.value.description),cssVars:n?void 0:d,themeClass:f?.themeClass,onRender:f?.onRender}},render(){let{$slots:e,mergedClsPrefix:t,onRender:n}=this;return n?.(),m(`div`,{class:[`${t}-empty`,this.themeClass],style:this.cssVars},this.showIcon?m(`div`,{class:`${t}-empty__icon`},e.icon?e.icon():m(A,{clsPrefix:t},{default:this.mergedRenderIcon})):null,this.showDescription?m(`div`,{class:`${t}-empty__description`},e.default?e.default():this.localizedDescription):null,e.extra?m(`div`,{class:`${t}-empty__extra`},e.extra()):null)}}),Ge={height:`calc(var(--n-option-height) * 7.6)`,paddingTiny:`4px 0`,paddingSmall:`4px 0`,paddingMedium:`4px 0`,paddingLarge:`4px 0`,paddingHuge:`4px 0`,optionPaddingTiny:`0 12px`,optionPaddingSmall:`0 12px`,optionPaddingMedium:`0 12px`,optionPaddingLarge:`0 12px`,optionPaddingHuge:`0 12px`,loadingSize:`18px`};function Ke(e){let{borderRadius:t,popoverColor:n,textColor3:r,dividerColor:i,textColor2:a,primaryColorPressed:o,textColorDisabled:s,primaryColor:c,opacityDisabled:l,hoverColor:u,fontSizeTiny:d,fontSizeSmall:f,fontSizeMedium:p,fontSizeLarge:m,fontSizeHuge:h,heightTiny:g,heightSmall:_,heightMedium:v,heightLarge:y,heightHuge:b}=e;return Object.assign(Object.assign({},Ge),{optionFontSizeTiny:d,optionFontSizeSmall:f,optionFontSizeMedium:p,optionFontSizeLarge:m,optionFontSizeHuge:h,optionHeightTiny:g,optionHeightSmall:_,optionHeightMedium:v,optionHeightLarge:y,optionHeightHuge:b,borderRadius:t,color:n,groupHeaderTextColor:r,actionDividerColor:i,optionTextColor:a,optionTextColorPressed:o,optionTextColorDisabled:s,optionTextColorActive:c,optionOpacityDisabled:l,optionCheckColor:c,optionColorPending:u,optionColorActive:`rgba(0, 0, 0, 0)`,optionColorActivePending:u,actionTextColor:a,loadingColor:c})}var qe=w({name:`InternalSelectMenu`,common:L,peers:{Scrollbar:de,Empty:He},self:Ke}),Je=f({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:n,labelFieldRef:r,nodePropsRef:i}=t(X);return{labelField:r,nodeProps:i,renderLabel:e,renderOption:n}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):$(i[this.labelField],i,!1),s=m(`div`,Object.assign({},a,{class:[`${e}-base-select-group-header`,a?.class]}),o);return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function Ye(e,t){return m(O,{name:`fade-in-scale-up-transition`},{default:()=>e?m(A,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>m(Le)}):null})}var Xe=f({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:n,pendingTmNodeRef:r,multipleRef:i,valueSetRef:a,renderLabelRef:o,renderOptionRef:s,labelFieldRef:c,valueFieldRef:l,showCheckmarkRef:u,nodePropsRef:d,handleOptionClick:f,handleOptionMouseEnter:p}=t(X),m=E(()=>{let{value:t}=r;return t?e.tmNode.key===t.key:!1});function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e;n.disabled||p(t,n)}function _(t){let{tmNode:n}=e,{value:r}=m;n.disabled||r||p(t,n)}return{multiple:i,isGrouped:E(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:u,nodeProps:d,isPending:m,isSelected:E(()=>{let{value:t}=n,{value:r}=i;if(t===null)return!1;let o=e.tmNode.rawNode[l.value];if(r){let{value:e}=a;return e.has(o)}else return t===o}),labelField:c,renderLabel:o,renderOption:s,handleMouseMove:_,handleMouseEnter:g,handleClick:h}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:a,nodeProps:o,renderOption:s,renderLabel:c,handleClick:l,handleMouseEnter:u,handleMouseMove:d}=this,f=Ye(n,e),p=c?[c(t,n),a&&f]:[$(t[this.labelField],t,n),a&&f],h=o?.(t),g=m(`div`,Object.assign({},h,{class:[`${e}-base-select-option`,t.class,h?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:a}],style:[h?.style||``,t.style||``],onClick:Ie([l,h?.onClick]),onMouseenter:Ie([u,h?.onMouseenter]),onMousemove:Ie([d,h?.onMousemove])}),m(`div`,{class:`${e}-base-select-option__content`},p));return t.render?t.render({node:g,option:t,selected:n}):s?s({node:g,option:t,selected:n}):g}}),Ze=V(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[V(`scrollbar`,`
 max-height: var(--n-height);
 `),V(`virtual-list`,`
 max-height: var(--n-height);
 `),V(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[r(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),V(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),V(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),r(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),r(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),r(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),r(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),V(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),V(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[B(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),v(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),v(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),B(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),B(`pending`,[v(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),B(`selected`,`
 color: var(--n-option-text-color-active);
 `,[v(`&::before`,`
 background-color: var(--n-option-color-active);
 `),B(`pending`,[v(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),B(`disabled`,`
 cursor: not-allowed;
 `,[l(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),B(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),r(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[xe({enterScale:`0.5`})])])]),Qe=f({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},S.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:r,mergedComponentPropsRef:o}=c(e),s=H(`InternalSelectMenu`,r,t),l=S(`InternalSelectMenu`,`-internal-select-menu`,Ze,qe,e,n(e,`clsPrefix`)),d=i(null),f=i(null),m=i(null),h=u(()=>e.treeMate.getFlattenedNodes()),g=u(()=>G(h.value)),_=i(null);function v(){let{treeMate:t}=e,n=null,{value:r}=e;r===null?n=t.getFirstAvailableNode():(n=e.multiple?t.getNode((r||[])[(r||[]).length-1]):t.getNode(r),(!n||n.disabled)&&(n=t.getFirstAvailableNode())),W(n||null)}function y(){let{value:t}=_;t&&!e.treeMate.getNode(t.key)&&(_.value=null)}let b;j(()=>e.show,t=>{t?b=j(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?v():y(),R(ne)):y()},{immediate:!0}):b?.()},{immediate:!0}),D(()=>{b?.()});let x=u(()=>_e(l.value.self[a(`optionHeight`,e.size)])),C=u(()=>Q(l.value.self[a(`padding`,e.size)])),w=u(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),T=u(()=>{let e=h.value;return e&&e.length===0}),E=u(()=>o?.value?.Select?.renderEmpty);function O(t){let{onToggle:n}=e;n&&n(t)}function A(t){let{onScroll:n}=e;n&&n(t)}function ee(e){var t;(t=m.value)==null||t.sync(),A(e)}function M(){var e;(e=m.value)==null||e.sync()}function P(){let{value:e}=_;return e||null}function F(e,t){t.disabled||W(t,!1)}function I(e,t){t.disabled||O(t)}function L(t){var n;U(t,`action`)||(n=e.onKeyup)==null||n.call(e,t)}function z(t){var n;U(t,`action`)||(n=e.onKeydown)==null||n.call(e,t)}function B(t){var n;(n=e.onMousedown)==null||n.call(e,t),!e.focusable&&t.preventDefault()}function V(){let{value:e}=_;e&&W(e.getNext({loop:!0}),!0)}function te(){let{value:e}=_;e&&W(e.getPrev({loop:!0}),!0)}function W(e,t=!1){_.value=e,t&&ne()}function ne(){var t,n;let r=_.value;if(!r)return;let i=g.value(r.key);i!==null&&(e.virtualScroll?(t=f.value)==null||t.scrollTo({index:i}):(n=m.value)==null||n.scrollTo({index:i,elSize:x.value}))}function K(t){var n;d.value?.contains(t.target)&&((n=e.onFocus)==null||n.call(e,t))}function q(t){var n;d.value?.contains(t.relatedTarget)||(n=e.onBlur)==null||n.call(e,t)}k(X,{handleOptionMouseEnter:F,handleOptionClick:I,valueSetRef:w,pendingTmNodeRef:_,nodePropsRef:n(e,`nodeProps`),showCheckmarkRef:n(e,`showCheckmark`),multipleRef:n(e,`multiple`),valueRef:n(e,`value`),renderLabelRef:n(e,`renderLabel`),renderOptionRef:n(e,`renderOption`),labelFieldRef:n(e,`labelField`),valueFieldRef:n(e,`valueField`)}),k(Y,d),N(()=>{let{value:e}=m;e&&e.sync()});let J=u(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:o,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:u,optionTextColor:d,optionTextColorDisabled:f,optionTextColorActive:p,optionOpacityDisabled:m,optionCheckColor:h,actionTextColor:g,optionColorPending:_,optionColorActive:v,loadingColor:y,loadingSize:b,optionColorActivePending:x,[a(`optionFontSize`,t)]:S,[a(`optionHeight`,t)]:C,[a(`optionPadding`,t)]:w}}=l.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":g,"--n-bezier":n,"--n-border-radius":i,"--n-color":o,"--n-option-font-size":S,"--n-group-header-text-color":s,"--n-option-check-color":h,"--n-option-color-pending":_,"--n-option-color-active":v,"--n-option-color-active-pending":x,"--n-option-height":C,"--n-option-opacity-disabled":m,"--n-option-text-color":d,"--n-option-text-color-active":p,"--n-option-text-color-disabled":f,"--n-option-text-color-pressed":u,"--n-option-padding":w,"--n-option-padding-left":Q(w,`left`),"--n-option-padding-right":Q(w,`right`),"--n-loading-color":y,"--n-loading-size":b}}),{inlineThemeDisabled:re}=e,ie=re?p(`internal-select-menu`,u(()=>e.size[0]),J,e):void 0,ae={selfRef:d,next:V,prev:te,getPendingTmNode:P};return Pe(d,e.onResize),Object.assign({mergedTheme:l,mergedClsPrefix:t,rtlEnabled:s,virtualListRef:f,scrollbarRef:m,itemSize:x,padding:C,flattenedNodes:h,empty:T,mergedRenderEmpty:E,virtualListContainer(){let{value:e}=f;return e?.listElRef},virtualListContent(){let{value:e}=f;return e?.itemsElRef},doScroll:A,handleFocusin:K,handleFocusout:q,handleKeyUp:L,handleKeyDown:z,handleMouseDown:B,handleVirtualListResize:M,handleVirtualListScroll:ee,cssVars:re?void 0:J,themeClass:ie?.themeClass,onRender:ie?.onRender},ae)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:r,themeClass:i,onRender:a}=this;return a?.(),m(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,i,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},_(e.header,e=>e&&m(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?m(`div`,{class:`${n}-base-select-menu__loading`},m(y,{clsPrefix:n,strokeWidth:20})):this.empty?m(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},I(e.empty,()=>[this.mergedRenderEmpty?.call(this)||m(We,{theme:r.peers.Empty,themeOverrides:r.peerOverrides.Empty,size:this.size})])):m(he,Object.assign({ref:`scrollbarRef`,theme:r.peers.Scrollbar,themeOverrides:r.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?m(Ne,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?m(Je,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:m(Xe,{clsPrefix:n,key:e.key,tmNode:e})}):m(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?m(Je,{key:e.key,clsPrefix:n,tmNode:e}):m(Xe,{clsPrefix:n,key:e.key,tmNode:e})))}),_(e.action,e=>e&&[m(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),m(ze,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),$e={closeIconSizeTiny:`12px`,closeIconSizeSmall:`12px`,closeIconSizeMedium:`14px`,closeIconSizeLarge:`14px`,closeSizeTiny:`16px`,closeSizeSmall:`16px`,closeSizeMedium:`18px`,closeSizeLarge:`18px`,padding:`0 7px`,closeMargin:`0 0 0 4px`};function et(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:c,baseColor:l,borderColor:u,opacityDisabled:d,tagColor:f,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,borderRadiusSmall:g,fontSizeMini:_,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,heightMini:x,heightTiny:S,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=e;return Object.assign(Object.assign({},$e),{closeBorderRadius:g,heightTiny:x,heightSmall:S,heightMedium:w,heightLarge:T,borderRadius:g,opacityDisabled:d,fontSizeTiny:_,fontSizeSmall:v,fontSizeMedium:y,fontSizeLarge:b,fontWeightStrong:A,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:l,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${u}`,textColor:t,color:f,colorBordered:`rgb(250, 250, 252)`,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${C(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:C(i,{alpha:.12}),colorBorderedPrimary:C(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:C(i,{alpha:.12}),closeColorPressedPrimary:C(i,{alpha:.18}),borderInfo:`1px solid ${C(a,{alpha:.3})}`,textColorInfo:a,colorInfo:C(a,{alpha:.12}),colorBorderedInfo:C(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:C(a,{alpha:.12}),closeColorPressedInfo:C(a,{alpha:.18}),borderSuccess:`1px solid ${C(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:C(o,{alpha:.12}),colorBorderedSuccess:C(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:C(o,{alpha:.12}),closeColorPressedSuccess:C(o,{alpha:.18}),borderWarning:`1px solid ${C(s,{alpha:.35})}`,textColorWarning:s,colorWarning:C(s,{alpha:.15}),colorBorderedWarning:C(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:C(s,{alpha:.12}),closeColorPressedWarning:C(s,{alpha:.18}),borderError:`1px solid ${C(c,{alpha:.23})}`,textColorError:c,colorError:C(c,{alpha:.1}),colorBorderedError:C(c,{alpha:.08}),closeIconColorError:c,closeIconColorHoverError:c,closeIconColorPressedError:c,closeColorHoverError:C(c,{alpha:.12}),closeColorPressedError:C(c,{alpha:.18})})}var tt={name:`Tag`,common:L,self:et},nt={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},rt=V(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[B(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),r(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),r(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),r(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),r(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),B(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[r(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),r(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),B(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),B(`icon, avatar`,[B(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),B(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),B(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[l(`disabled`,[v(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[l(`checked`,`color: var(--n-text-color-hover-checkable);`)]),v(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[l(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),B(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[l(`disabled`,[v(`&:hover`,`background-color: var(--n-color-checked-hover);`),v(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),it=Object.assign(Object.assign(Object.assign({},S.props),nt),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),at=F(`n-tag`),ot=f({name:`Tag`,props:it,slots:Object,setup(e){let t=i(null),{mergedBorderedRef:r,mergedClsPrefixRef:o,inlineThemeDisabled:s,mergedRtlRef:l,mergedComponentPropsRef:d}=c(e),f=u(()=>e.size||d?.value?.Tag?.size||`medium`),m=S(`Tag`,`-tag`,rt,tt,e,o);k(at,{roundRef:n(e,`round`)});function h(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function g(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&P(n,t)}}let _={setTextContent(e){let{value:n}=t;n&&(n.textContent=e)}},v=H(`Tag`,l,o),y=u(()=>{let{type:t,color:{color:n,textColor:i}={}}=e,o=f.value,{common:{cubicBezierEaseInOut:s},self:{padding:c,closeMargin:l,borderRadius:u,opacityDisabled:d,textColorCheckable:p,textColorHoverCheckable:h,textColorPressedCheckable:g,textColorChecked:_,colorCheckable:v,colorHoverCheckable:y,colorPressedCheckable:b,colorChecked:x,colorCheckedHover:S,colorCheckedPressed:C,closeBorderRadius:w,fontWeightStrong:T,[a(`colorBordered`,t)]:E,[a(`closeSize`,o)]:D,[a(`closeIconSize`,o)]:O,[a(`fontSize`,o)]:k,[a(`height`,o)]:A,[a(`color`,t)]:j,[a(`textColor`,t)]:ee,[a(`border`,t)]:M,[a(`closeIconColor`,t)]:N,[a(`closeIconColorHover`,t)]:P,[a(`closeIconColorPressed`,t)]:F,[a(`closeColorHover`,t)]:I,[a(`closeColorPressed`,t)]:L}}=m.value,R=Q(l);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${A} - 8px)`,"--n-bezier":s,"--n-border-radius":u,"--n-border":M,"--n-close-icon-size":O,"--n-close-color-pressed":L,"--n-close-color-hover":I,"--n-close-border-radius":w,"--n-close-icon-color":N,"--n-close-icon-color-hover":P,"--n-close-icon-color-pressed":F,"--n-close-icon-color-disabled":N,"--n-close-margin-top":R.top,"--n-close-margin-right":R.right,"--n-close-margin-bottom":R.bottom,"--n-close-margin-left":R.left,"--n-close-size":D,"--n-color":n||(r.value?E:j),"--n-color-checkable":v,"--n-color-checked":x,"--n-color-checked-hover":S,"--n-color-checked-pressed":C,"--n-color-hover-checkable":y,"--n-color-pressed-checkable":b,"--n-font-size":k,"--n-height":A,"--n-opacity-disabled":d,"--n-padding":c,"--n-text-color":i||ee,"--n-text-color-checkable":p,"--n-text-color-checked":_,"--n-text-color-hover-checkable":h,"--n-text-color-pressed-checkable":g}}),b=s?p(`tag`,u(()=>{let t=``,{type:n,color:{color:i,textColor:a}={}}=e;return t+=n[0],t+=f.value[0],i&&(t+=`a${z(i)}`),a&&(t+=`b${z(a)}`),r.value&&(t+=`c`),t}),y,e):void 0;return Object.assign(Object.assign({},_),{rtlEnabled:v,mergedClsPrefix:o,contentRef:t,mergedBordered:r,handleClick:h,handleCloseClick:g,cssVars:s?void 0:y,themeClass:b?.themeClass,onRender:b?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:a,onRender:o,$slots:s}=this;o?.();let c=_(s.avatar,e=>e&&m(`div`,{class:`${t}-tag__avatar`},e)),l=_(s.icon,e=>e&&m(`div`,{class:`${t}-tag__icon`},e));return m(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:a,[`${t}-tag--avatar`]:c,[`${t}-tag--icon`]:l,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},l||c,m(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?m(me,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:a,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?m(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),st={paddingSingle:`0 26px 0 12px`,paddingMultiple:`3px 26px 0 12px`,clearSize:`16px`,arrowSize:`16px`};function ct(e){let{borderRadius:t,textColor2:n,textColorDisabled:r,inputColor:i,inputColorDisabled:a,primaryColor:o,primaryColorHover:s,warningColor:c,warningColorHover:l,errorColor:u,errorColorHover:d,borderColor:f,iconColor:p,iconColorDisabled:m,clearColor:h,clearColorHover:g,clearColorPressed:_,placeholderColor:v,placeholderColorDisabled:y,fontSizeTiny:b,fontSizeSmall:x,fontSizeMedium:S,fontSizeLarge:w,heightTiny:T,heightSmall:E,heightMedium:D,heightLarge:O,fontWeight:k}=e;return Object.assign(Object.assign({},st),{fontSizeTiny:b,fontSizeSmall:x,fontSizeMedium:S,fontSizeLarge:w,heightTiny:T,heightSmall:E,heightMedium:D,heightLarge:O,borderRadius:t,fontWeight:k,textColor:n,textColorDisabled:r,placeholderColor:v,placeholderColorDisabled:y,color:i,colorDisabled:a,colorActive:i,border:`1px solid ${f}`,borderHover:`1px solid ${s}`,borderActive:`1px solid ${o}`,borderFocus:`1px solid ${s}`,boxShadowHover:`none`,boxShadowActive:`0 0 0 2px ${C(o,{alpha:.2})}`,boxShadowFocus:`0 0 0 2px ${C(o,{alpha:.2})}`,caretColor:o,arrowColor:p,arrowColorDisabled:m,loadingColor:o,borderWarning:`1px solid ${c}`,borderHoverWarning:`1px solid ${l}`,borderActiveWarning:`1px solid ${c}`,borderFocusWarning:`1px solid ${l}`,boxShadowHoverWarning:`none`,boxShadowActiveWarning:`0 0 0 2px ${C(c,{alpha:.2})}`,boxShadowFocusWarning:`0 0 0 2px ${C(c,{alpha:.2})}`,colorActiveWarning:i,caretColorWarning:c,borderError:`1px solid ${u}`,borderHoverError:`1px solid ${d}`,borderActiveError:`1px solid ${u}`,borderFocusError:`1px solid ${d}`,boxShadowHoverError:`none`,boxShadowActiveError:`0 0 0 2px ${C(u,{alpha:.2})}`,boxShadowFocusError:`0 0 0 2px ${C(u,{alpha:.2})}`,colorActiveError:i,caretColorError:u,clearColor:h,clearColorHover:g,clearColorPressed:_})}var lt=w({name:`InternalSelection`,common:L,peers:{Popover:ae},self:ct}),ut=v([V(`base-selection`,`
 --n-padding-single: var(--n-padding-single-top) var(--n-padding-single-right) var(--n-padding-single-bottom) var(--n-padding-single-left);
 --n-padding-multiple: var(--n-padding-multiple-top) var(--n-padding-multiple-right) var(--n-padding-multiple-bottom) var(--n-padding-multiple-left);
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[V(`base-loading`,`
 color: var(--n-loading-color);
 `),V(`base-selection-tags`,`min-height: var(--n-height);`),r(`border, state-border`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),r(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),V(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[r(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),V(`base-selection-overlay`,`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[r(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),V(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[r(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),V(`base-selection-tags`,`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),V(`base-selection-label`,`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[V(`base-selection-input`,`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[r(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),r(`render-label`,`
 color: var(--n-text-color);
 `)]),l(`disabled`,[v(`&:hover`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),B(`focus`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),B(`active`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),V(`base-selection-label`,`background-color: var(--n-color-active);`),V(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),B(`disabled`,`cursor: not-allowed;`,[r(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),V(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[V(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),r(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),V(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),V(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),V(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[r(`input`,`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),r(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>B(`${e}-status`,[r(`state-border`,`border: var(--n-border-${e});`),l(`disabled`,[v(`&:hover`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),B(`active`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),V(`base-selection-label`,`background-color: var(--n-color-active-${e});`),V(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),B(`focus`,[r(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),V(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),V(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[v(`&:last-child`,`padding-right: 0;`),V(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[r(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),dt=f({name:`InternalSelection`,props:Object.assign(Object.assign({},S.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:r}=c(e),o=H(`InternalSelection`,r,t),s=i(null),l=i(null),d=i(null),f=i(null),m=i(null),h=i(null),g=i(null),_=i(null),v=i(null),y=i(null),b=i(!1),x=i(!1),C=i(!1),w=S(`InternalSelection`,`-internal-selection`,ut,lt,e,n(e,`clsPrefix`)),E=u(()=>e.clearable&&!e.disabled&&(C.value||e.active)),D=u(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):$(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),O=u(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),k=u(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function A(){var t;let{value:n}=s;if(n){let{value:r}=l;r&&(r.style.width=`${n.offsetWidth}px`,e.maxTagCount!==`responsive`&&((t=v.value)==null||t.sync({showAllItemsBeforeCalculate:!1})))}}function ee(){let{value:e}=y;e&&(e.style.display=`none`)}function M(){let{value:e}=y;e&&(e.style.display=`inline-block`)}j(n(e,`active`),e=>{e||ee()}),j(n(e,`pattern`),()=>{e.multiple&&R(A)});function P(t){let{onFocus:n}=e;n&&n(t)}function F(t){let{onBlur:n}=e;n&&n(t)}function I(t){let{onDeleteOption:n}=e;n&&n(t)}function L(t){let{onClear:n}=e;n&&n(t)}function z(t){let{onPatternInput:n}=e;n&&n(t)}function B(e){(!e.relatedTarget||!d.value?.contains(e.relatedTarget))&&P(e)}function V(e){d.value?.contains(e.relatedTarget)||F(e)}function te(e){L(e)}function U(){C.value=!0}function W(){C.value=!1}function ne(t){!e.active||!e.filterable||t.target!==l.value&&t.preventDefault()}function G(e){I(e)}let K=i(!1);function q(t){if(t.key===`Backspace`&&!K.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&G(t[t.length-1])}}let J=null;function re(t){let{value:n}=s;n&&(n.textContent=t.target.value,A()),e.ignoreComposition&&K.value?J=t:z(t)}function ie(){K.value=!0}function ae(){K.value=!1,e.ignoreComposition&&z(J),J=null}function oe(t){var n;x.value=!0,(n=e.onPatternFocus)==null||n.call(e,t)}function se(t){var n;x.value=!1,(n=e.onPatternBlur)==null||n.call(e,t)}function ce(){var t,n;if(e.filterable)x.value=!1,(t=h.value)==null||t.blur(),(n=l.value)==null||n.blur();else if(e.multiple){let{value:e}=f;e?.blur()}else{let{value:e}=m;e?.blur()}}function le(){var t,n,r;e.filterable?(x.value=!1,(t=h.value)==null||t.focus()):e.multiple?(n=f.value)==null||n.focus():(r=m.value)==null||r.focus()}function Y(){let{value:e}=l;e&&(M(),e.focus())}function X(){let{value:e}=l;e&&e.blur()}function Z(e){let{value:t}=g;t&&t.setTextContent(`+${e}`)}function ue(){let{value:e}=_;return e}function de(){return l.value}let fe=null;function pe(){fe!==null&&window.clearTimeout(fe)}function me(){e.active||(pe(),fe=window.setTimeout(()=>{k.value&&(b.value=!0)},100))}function he(){pe()}function ge(e){e||(pe(),b.value=!1)}j(k,e=>{e||(b.value=!1)}),N(()=>{T(()=>{let t=h.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=x.value?-1:0)})}),Pe(d,e.onResize);let{inlineThemeDisabled:_e}=e,ve=u(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:o,placeholderColor:s,textColor:c,paddingSingle:l,paddingMultiple:u,caretColor:d,colorDisabled:f,textColorDisabled:p,placeholderColorDisabled:m,colorActive:h,boxShadowFocus:g,boxShadowActive:_,boxShadowHover:v,border:y,borderFocus:b,borderHover:x,borderActive:S,arrowColor:C,arrowColorDisabled:T,loadingColor:E,colorActiveWarning:D,boxShadowFocusWarning:O,boxShadowActiveWarning:k,boxShadowHoverWarning:A,borderWarning:j,borderFocusWarning:ee,borderHoverWarning:M,borderActiveWarning:N,colorActiveError:P,boxShadowFocusError:F,boxShadowActiveError:I,boxShadowHoverError:L,borderError:R,borderFocusError:z,borderHoverError:B,borderActiveError:V,clearColor:H,clearColorHover:te,clearColorPressed:U,clearSize:W,arrowSize:ne,[a(`height`,t)]:G,[a(`fontSize`,t)]:K}}=w.value,q=Q(l),J=Q(u);return{"--n-bezier":n,"--n-border":y,"--n-border-active":S,"--n-border-focus":b,"--n-border-hover":x,"--n-border-radius":i,"--n-box-shadow-active":_,"--n-box-shadow-focus":g,"--n-box-shadow-hover":v,"--n-caret-color":d,"--n-color":o,"--n-color-active":h,"--n-color-disabled":f,"--n-font-size":K,"--n-height":G,"--n-padding-single-top":q.top,"--n-padding-multiple-top":J.top,"--n-padding-single-right":q.right,"--n-padding-multiple-right":J.right,"--n-padding-single-left":q.left,"--n-padding-multiple-left":J.left,"--n-padding-single-bottom":q.bottom,"--n-padding-multiple-bottom":J.bottom,"--n-placeholder-color":s,"--n-placeholder-color-disabled":m,"--n-text-color":c,"--n-text-color-disabled":p,"--n-arrow-color":C,"--n-arrow-color-disabled":T,"--n-loading-color":E,"--n-color-active-warning":D,"--n-box-shadow-focus-warning":O,"--n-box-shadow-active-warning":k,"--n-box-shadow-hover-warning":A,"--n-border-warning":j,"--n-border-focus-warning":ee,"--n-border-hover-warning":M,"--n-border-active-warning":N,"--n-color-active-error":P,"--n-box-shadow-focus-error":F,"--n-box-shadow-active-error":I,"--n-box-shadow-hover-error":L,"--n-border-error":R,"--n-border-focus-error":z,"--n-border-hover-error":B,"--n-border-active-error":V,"--n-clear-size":W,"--n-clear-color":H,"--n-clear-color-hover":te,"--n-clear-color-pressed":U,"--n-arrow-size":ne,"--n-font-weight":r}}),ye=_e?p(`internal-selection`,u(()=>e.size[0]),ve,e):void 0;return{mergedTheme:w,mergedClearable:E,mergedClsPrefix:t,rtlEnabled:o,patternInputFocused:x,filterablePlaceholder:D,label:O,selected:k,showTagsPanel:b,isComposing:K,counterRef:g,counterWrapperRef:_,patternInputMirrorRef:s,patternInputRef:l,selfRef:d,multipleElRef:f,singleElRef:m,patternInputWrapperRef:h,overflowRef:v,inputTagElRef:y,handleMouseDown:ne,handleFocusin:B,handleClear:te,handleMouseEnter:U,handleMouseLeave:W,handleDeleteOption:G,handlePatternKeyDown:q,handlePatternInputInput:re,handlePatternInputBlur:se,handlePatternInputFocus:oe,handleMouseEnterCounter:me,handleMouseLeaveCounter:he,handleFocusout:V,handleCompositionEnd:ae,handleCompositionStart:ie,onPopoverUpdateShow:ge,focus:le,focusInput:Y,blur:ce,blurInput:X,updateCounter:Z,getCounter:ue,getTail:de,renderLabel:e.renderLabel,cssVars:_e?void 0:ve,themeClass:ye?.themeClass,onRender:ye?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:i,maxTagCount:a,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:u,renderTag:d,renderLabel:f}=this;u?.();let p=a===`responsive`,h=typeof a==`number`,g=p||h,_=m(fe,null,{default:()=>m(ye,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),v;if(t){let{labelField:e}=this,t=t=>m(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},d?d({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):m(ot,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>f?f(t,!0):$(t[e],t,!0)})),s=()=>(h?this.selectedOptions.slice(0,a):this.selectedOptions).map(t),u=i?m(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},m(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),m(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,y=p?()=>m(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},m(ot,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,b;if(h){let e=this.selectedOptions.length-a;e>0&&(b=m(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},m(ot,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let x=p?i?m(K,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y,tail:()=>u}):m(K,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y}):h&&b?s().concat(b):s(),S=g?()=>m(`div`,{class:`${c}-base-selection-popover`},p?s():this.selectedOptions.map(t)):void 0,C=g?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,w=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?m(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},m(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,T=i?m(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},x,p?null:u,_):m(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},x,_);v=m(o,null,g?m(se,Object.assign({},C,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>T,default:S}):T,w)}else if(i){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;v=m(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:le(this.label)},m(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?m(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},m(`div`,{class:`${c}-base-selection-overlay__wrapper`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):$(this.label,this.selectedOption,!0))):null,t?m(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},m(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,_)}else v=m(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?m(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},m(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):m(`div`,{class:`${c}-base-selection-input`,title:le(this.label),key:`input`},m(`div`,{class:`${c}-base-selection-input__content`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):$(this.label,this.selectedOption,!0))),_);return m(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},v,s?m(`div`,{class:`${c}-base-selection__border`}):null,s?m(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ft(e){return e.type===`group`}function pt(e){return e.type===`ignored`}function mt(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function ht(e,t){return{getIsGroup:ft,getIgnored:pt,getKey(t){return ft(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function gt(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ft(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(pt(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function _t(e,t,n){let r=new Map;return e.forEach(e=>{ft(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}function vt(e){let{boxShadow2:t}=e;return{menuBoxShadow:t}}var yt=w({name:`Select`,common:L,peers:{InternalSelection:lt,InternalSelectMenu:qe},self:vt}),bt=v([V(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),V(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[xe({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),xt=f({name:`Select`,props:Object.assign(Object.assign({},S.props),{to:W.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(t){let{mergedClsPrefixRef:r,mergedBorderedRef:a,namespaceRef:o,inlineThemeDisabled:s,mergedComponentPropsRef:l}=c(t),f=S(`Select`,`-select`,bt,yt,t,r),m=i(t.defaultValue),h=e(n(t,`value`),m),g=i(!1),_=i(``),v=ne(t,[`items`,`options`]),y=i([]),b=i([]),x=u(()=>b.value.concat(y.value).concat(v.value)),C=u(()=>{let{filter:e}=t;if(e)return e;let{labelField:n,valueField:r}=t;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return mt(e,i);let a=t[r];return typeof a==`string`?mt(e,a):typeof a==`number`?mt(e,String(a)):!1}}),w=u(()=>{if(t.remote)return v.value;{let{value:e}=x,{value:n}=_;return!n.length||!t.filterable?e:gt(e,C.value,n,t.childrenField)}}),T=u(()=>{let{valueField:e,childrenField:n}=t,r=ht(e,n);return ce(w.value,r)}),E=u(()=>_t(x.value,t.valueField,t.childrenField)),D=i(!1),O=e(n(t,`show`),D),k=i(null),A=i(null),M=i(null),{localeRef:N}=ve(`Select`),F=u(()=>t.placeholder??N.value.placeholder),I=[],L=i(new Map),R=u(()=>{let{fallbackOption:e}=t;if(e===void 0){let{labelField:e,valueField:n}=t;return t=>({[e]:String(t),[n]:t})}return e===!1?!1:t=>Object.assign(e(t),{value:t})});function z(e){let n=t.remote,{value:r}=L,{value:i}=E,{value:a}=R,o=[];return e.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let B=u(()=>{if(t.multiple){let{value:e}=h;return Array.isArray(e)?z(e):[]}return null}),V=u(()=>{let{value:e}=h;return!t.multiple&&!Array.isArray(e)?e===null?null:z([e])[0]||null:null}),H=d(t,{mergedSize:e=>{let{size:n}=t;if(n)return n;let{mergedSize:r}=e||{};return r?.value?r.value:l?.value?.Select?.size||`medium`}}),{mergedSizeRef:te,mergedDisabledRef:G,mergedStatusRef:K}=H;function q(e,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=t,{nTriggerFormChange:o,nTriggerFormInput:s}=H;r&&P(r,e,n),a&&P(a,e,n),i&&P(i,e,n),m.value=e,o(),s()}function J(e){let{onBlur:n}=t,{nTriggerFormBlur:r}=H;n&&P(n,e),r()}function re(){let{onClear:e}=t;e&&P(e)}function ie(e){let{onFocus:n,showOnFocus:r}=t,{nTriggerFormFocus:i}=H;n&&P(n,e),i(),r&&Y()}function ae(e){let{onSearch:n}=t;n&&P(n,e)}function oe(e){let{onScroll:n}=t;n&&P(n,e)}function se(){var e;let{remote:n,multiple:r}=t;if(n){let{value:n}=L;if(r){let{valueField:r}=t;(e=B.value)==null||e.forEach(e=>{n.set(e[r],e)})}else{let e=V.value;e&&n.set(e[t.valueField],e)}}}function le(e){let{onUpdateShow:n,"onUpdate:show":r}=t;n&&P(n,e),r&&P(r,e),D.value=e}function Y(){G.value||(le(!0),D.value=!0,t.filterable&&je())}function X(){le(!1)}function Z(){_.value=``,b.value=I}let Q=i(!1);function ue(){t.filterable&&(Q.value=!0)}function de(){t.filterable&&(Q.value=!1,O.value||Z())}function fe(){G.value||(O.value?t.filterable?je():X():Y())}function me(e){(M.value?.selfRef)?.contains(e.relatedTarget)||(g.value=!1,J(e),X())}function he(e){ie(e),g.value=!0}function ge(){g.value=!0}function _e(e){k.value?.$el.contains(e.relatedTarget)||(g.value=!1,J(e),X())}function ye(){var e;(e=k.value)==null||e.focus(),X()}function be(e){O.value&&(k.value?.$el.contains(pe(e))||X())}function $(e){if(!Array.isArray(e))return[];if(R.value)return Array.from(e);{let{remote:n}=t,{value:r}=E;if(n){let{value:t}=L;return e.filter(e=>r.has(e)||t.has(e))}else return e.filter(e=>r.has(e))}}function xe(e){Ce(e.rawNode)}function Ce(e){if(G.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=t;if(n&&!r){let{value:e}=b,t=e[0]||null;if(t){let e=y.value;e.length?e.push(t):y.value=[t],b.value=I}}if(r&&L.value.set(e[a],e),t.multiple){let t=$(h.value),o=t.findIndex(t=>t===e[a]);if(~o){if(t.splice(o,1),n&&!r){let t=we(e[a]);~t&&(y.value.splice(t,1),i&&(_.value=``))}}else t.push(e[a]),i&&(_.value=``);q(t,z(t))}else{if(n&&!r){let t=we(e[a]);~t?y.value=[y.value[t]]:y.value=I}Ae(),X(),q(e[a],e)}}function we(e){return y.value.findIndex(n=>n[t.valueField]===e)}function Te(e){O.value||Y();let{value:n}=e.target;_.value=n;let{tag:r,remote:i}=t;if(ae(n),r&&!i){if(!n){b.value=I;return}let{onCreate:e}=t,r=e?e(n):{[t.labelField]:n,[t.valueField]:n},{valueField:i,labelField:a}=t;v.value.some(e=>e[i]===r[i]||e[a]===r[a])||y.value.some(e=>e[i]===r[i]||e[a]===r[a])?b.value=I:b.value=[r]}}function Ee(e){e.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=t;!n&&t.filterable&&X(),r&&!i&&a&&(y.value=I),re(),n?q([],[]):q(null,null)}function De(e){!U(e,`action`)&&!U(e,`empty`)&&!U(e,`header`)&&e.preventDefault()}function Oe(e){oe(e)}function ke(e){var n,r,i;if(!t.keyboard){e.preventDefault();return}switch(e.key){case` `:if(t.filterable)break;e.preventDefault();case`Enter`:if(!k.value?.isComposing){if(O.value){let e=M.value?.getPendingTmNode();e?xe(e):t.filterable||(X(),Ae())}else if(Y(),t.tag&&Q.value){let e=b.value[0];if(e){let n=e[t.valueField],{value:r}=h;t.multiple&&Array.isArray(r)&&r.includes(n)||Ce(e)}}}e.preventDefault();break;case`ArrowUp`:if(e.preventDefault(),t.loading)return;O.value&&((n=M.value)==null||n.prev());break;case`ArrowDown`:if(e.preventDefault(),t.loading)return;O.value?(r=M.value)==null||r.next():Y();break;case`Escape`:O.value&&(Se(e),X()),(i=k.value)==null||i.focus();break}}function Ae(){var e;(e=k.value)==null||e.focus()}function je(){var e;(e=k.value)==null||e.focusInput()}function Me(){var e;O.value&&((e=A.value)==null||e.syncPosition())}se(),j(n(t,`options`),se);let Ne={focus:()=>{var e;(e=k.value)==null||e.focus()},focusInput:()=>{var e;(e=k.value)==null||e.focusInput()},blur:()=>{var e;(e=k.value)==null||e.blur()},blurInput:()=>{var e;(e=k.value)==null||e.blurInput()}},Pe=u(()=>{let{self:{menuBoxShadow:e}}=f.value;return{"--n-menu-box-shadow":e}}),Fe=s?p(`select`,void 0,Pe,t):void 0;return Object.assign(Object.assign({},Ne),{mergedStatus:K,mergedClsPrefix:r,mergedBordered:a,namespace:o,treeMate:T,isMounted:ee(),triggerRef:k,menuRef:M,pattern:_,uncontrolledShow:D,mergedShow:O,adjustedTo:W(t),uncontrolledValue:m,mergedValue:h,followerRef:A,localizedPlaceholder:F,selectedOption:V,selectedOptions:B,mergedSize:te,mergedDisabled:G,focused:g,activeWithoutMenuOpen:Q,inlineThemeDisabled:s,onTriggerInputFocus:ue,onTriggerInputBlur:de,handleTriggerOrMenuResize:Me,handleMenuFocus:ge,handleMenuBlur:_e,handleMenuTabOut:ye,handleTriggerClick:fe,handleToggle:xe,handleDeleteOption:Ce,handlePatternInput:Te,handleClear:Ee,handleTriggerBlur:me,handleTriggerFocus:he,handleKeydown:ke,handleMenuAfterLeave:Z,handleMenuClickOutside:be,handleMenuScroll:Oe,handleMenuKeydown:ke,handleMenuMousedown:De,mergedTheme:f,cssVars:s?void 0:Pe,themeClass:Fe?.themeClass,onRender:Fe?.onRender})},render(){return m(`div`,{class:`${this.mergedClsPrefix}-select`},m(J,null,{default:()=>[m(re,null,{default:()=>m(dt,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),m(q,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===W.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>m(O,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),h(m(Qe,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[s,this.mergedShow],[be,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[be,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{Qe as a,He as c,Fe as d,Ne as f,ot as i,ze as l,yt as n,qe as o,ht as r,We as s,xt as t,Ie as u};