import{E as e,I as t,L as n,O as r,P as i,S as a,T as o,h as s,p as c,z as l}from"./token-DvMyJbbF.js";import{$ as u,Bn as d,Cn as f,Dn as p,Gt as m,Ht as h,Kt as g,Ln as _,Mn as v,On as y,Q as b,Qt as x,Rn as S,Sn as C,Tn as w,Vt as T,Wt as E,Xn as D,_ as O,_t as k,at as A,b as j,bn as M,bt as N,cn as P,er as F,et as I,g as L,gn as R,h as z,ht as ee,in as B,l as V,lt as H,mt as te,qt as U,st as W,tn as G,u as K,vn as q,vt as J,xn as Y,yt as X}from"./Button-scR9_EKp.js";import{A as ne,B as re,D as Z,H as ie,J as Q,M as $,N as ae,O as oe,S as se,T as ce,V as le,_ as ue,c as de,i as fe,j as pe,l as me,p as he,q as ge,u as _e,w as ve}from"./Tooltip-DPODkf4E.js";import{l as ye,m as be,u as xe}from"./Card-Ds-uU5ZB.js";function Se(e){return e&-e}var Ce=class{constructor(e,t){this.l=e,this.min=t;let n=Array(e+1);for(let t=0;t<e+1;++t)n[t]=0;this.ft=n}add(e,t){if(t===0)return;let{l:n,ft:r}=this;for(e+=1;e<=n;)r[e]+=t,e+=Se(e)}get(e){return this.sum(e+1)-this.sum(e)}sum(e){if(e===void 0&&(e=this.l),e<=0)return 0;let{ft:t,min:n,l:r}=this;if(e>r)throw Error("[FinweckTree.sum]: `i` is larger than length.");let i=e*n;for(;e>0;)i+=t[e],e-=Se(e);return i}getBound(e){let t=0,n=this.l;for(;n>t;){let r=Math.floor((t+n)/2),i=this.sum(r);if(i>e){n=r;continue}else if(i<e){if(t===r)return this.sum(t+1)<=e?t+1:r;t=r}else return r}return t}},we;function Te(){return typeof document>`u`?!1:(we===void 0&&(we=`matchMedia`in window?window.matchMedia(`(pointer:coarse)`).matches:!1),we)}var Ee;function De(){return typeof document>`u`?1:(Ee===void 0&&(Ee=`chrome`in window?window.devicePixelRatio:1),Ee)}var Oe=`VVirtualListXScroll`;function ke({columnsRef:e,renderColRef:t,renderItemWithColsRef:n}){let r=D(0),i=D(0),a=P(()=>{let t=e.value;if(t.length===0)return null;let n=new Ce(t.length,0);return t.forEach((e,t)=>{n.add(t,e.width)}),n});return v(Oe,{startIndexRef:X(()=>{let e=a.value;return e===null?0:Math.max(e.getBound(i.value)-1,0)}),endIndexRef:X(()=>{let t=a.value;return t===null?0:Math.min(t.getBound(i.value+r.value)+1,e.value.length-1)}),columnsRef:e,renderColRef:t,renderItemWithColsRef:n,getLeft:e=>{let t=a.value;return t===null?0:t.sum(e)}}),{listWidthRef:r,scrollLeftRef:i}}var Ae=R({name:`VirtualListRow`,props:{index:{type:Number,required:!0},item:{type:Object,required:!0}},setup(){let{startIndexRef:e,endIndexRef:t,columnsRef:n,getLeft:r,renderColRef:i,renderItemWithColsRef:a}=M(Oe);return{startIndex:e,endIndex:t,columns:n,renderCol:i,renderItemWithCols:a,getLeft:r}},render(){let{startIndex:e,endIndex:t,columns:n,renderCol:r,renderItemWithCols:i,getLeft:a,item:o}=this;if(i!=null)return i({itemIndex:this.index,startColIndex:e,endColIndex:t,allColumns:n,item:o,getLeft:a});if(r!=null){let i=[];for(let s=e;s<=t;++s){let e=n[s];i.push(r({column:e,left:a(s),item:o}))}return i}return null}}),je=Z(`.v-vl`,{maxHeight:`inherit`,height:`100%`,overflow:`auto`,minWidth:`1px`},[Z(`&:not(.v-vl--show-scrollbar)`,{scrollbarWidth:`none`},[Z(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,{width:0,height:0,display:`none`})])]),Me=R({name:`VirtualList`,inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},columns:{type:Array,default:()=>[]},renderCol:Function,renderItemWithCols:Function,items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:`div`},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:`key`},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(e){let t=ee();je.mount({id:`vueuc/virtual-list`,head:!0,anchorMetaName:oe,ssr:t}),y(()=>{let{defaultScrollIndex:t,defaultScrollKey:n}=e;t==null?n!=null&&x({key:n}):x({index:t})});let r=!1,a=!1;f(()=>{if(r=!1,!a){a=!0;return}x({top:_.value,left:c.value})}),p(()=>{r=!0,a||=!0});let o=X(()=>{if(e.renderCol==null&&e.renderItemWithCols==null||e.columns.length===0)return;let t=0;return e.columns.forEach(e=>{t+=e.width}),t}),s=P(()=>{let t=new Map,{keyField:n}=e;return e.items.forEach((e,r)=>{t.set(e[n],r)}),t}),{scrollLeftRef:c,listWidthRef:l}=ke({columnsRef:F(e,`columns`),renderColRef:F(e,`renderCol`),renderItemWithColsRef:F(e,`renderItemWithCols`)}),u=D(null),d=D(void 0),m=new Map,h=P(()=>{let{items:t,itemSize:n,keyField:r}=e,i=new Ce(t.length,n);return t.forEach((e,t)=>{let n=e[r],a=m.get(n);a!==void 0&&i.add(t,a)}),i}),g=D(0),_=D(0),v=X(()=>Math.max(h.value.getBound(_.value-i(e.paddingTop))-1,0)),b=P(()=>{let{value:t}=d;if(t===void 0)return[];let{items:n,itemSize:r}=e,i=v.value,a=Math.min(i+Math.ceil(t/r+1),n.length-1),o=[];for(let e=i;e<=a;++e)o.push(n[e]);return o}),x=(e,t)=>{if(typeof e==`number`){T(e,t,`auto`);return}let{left:n,top:r,index:i,key:a,position:o,behavior:c,debounce:l=!0}=e;if(n!==void 0||r!==void 0)T(n,r,c);else if(i!==void 0)w(i,c,l);else if(a!==void 0){let e=s.value.get(a);e!==void 0&&w(e,c,l)}else o===`bottom`?T(0,2**53-1,c):o===`top`&&T(0,0,c)},S,C=null;function w(t,n,r){let{value:a}=h,o=a.sum(t)+i(e.paddingTop);if(!r)u.value.scrollTo({left:0,top:o,behavior:n});else{S=t,C!==null&&window.clearTimeout(C),C=window.setTimeout(()=>{S=void 0,C=null},16);let{scrollTop:e,offsetHeight:r}=u.value;if(o>e){let i=a.get(t);o+i<=e+r||u.value.scrollTo({left:0,top:o+i-r,behavior:n})}else u.value.scrollTo({left:0,top:o,behavior:n})}}function T(e,t,n){u.value.scrollTo({left:e,top:t,behavior:n})}function E(t,n){if(r||e.ignoreItemResize||I(n.target))return;let{value:i}=h,a=s.value.get(t),o=i.get(a),c=n.borderBoxSize?.[0]?.blockSize??n.contentRect.height;if(c===o)return;c-e.itemSize===0?m.delete(t):m.set(t,c-e.itemSize);let l=c-o;if(l===0)return;i.add(a,l);let d=u.value;if(d!=null){if(S===void 0){let e=i.sum(a);d.scrollTop>e&&d.scrollBy(0,l)}else (a<S||a===S&&c+i.sum(a)>d.scrollTop+d.offsetHeight)&&d.scrollBy(0,l);N()}g.value++}let O=!Te(),k=!1;function A(t){var n;(n=e.onScroll)==null||n.call(e,t),(!O||!k)&&N()}function j(t){var n;if((n=e.onWheel)==null||n.call(e,t),O){let e=u.value;if(e!=null){if(t.deltaX===0&&(e.scrollTop===0&&t.deltaY<=0||e.scrollTop+e.offsetHeight>=e.scrollHeight&&t.deltaY>=0))return;t.preventDefault(),e.scrollTop+=t.deltaY/De(),e.scrollLeft+=t.deltaX/De(),N(),k=!0,Q(()=>{k=!1})}}}function M(t){if(r||I(t.target))return;if(e.renderCol==null&&e.renderItemWithCols==null){if(t.contentRect.height===d.value)return}else if(t.contentRect.height===d.value&&t.contentRect.width===l.value)return;d.value=t.contentRect.height,l.value=t.contentRect.width;let{onResize:n}=e;n!==void 0&&n(t)}function N(){let{value:e}=u;e!=null&&(_.value=e.scrollTop,c.value=e.scrollLeft)}function I(e){let t=e;for(;t!==null;){if(t.style.display===`none`)return!0;t=t.parentElement}return!1}return{listHeight:d,listStyle:{overflow:`auto`},keyToIndex:s,itemsStyle:P(()=>{let{itemResizable:t}=e,r=n(h.value.sum());return g.value,[e.itemsStyle,{boxSizing:`content-box`,width:n(o.value),height:t?``:r,minHeight:t?r:``,paddingTop:n(e.paddingTop),paddingBottom:n(e.paddingBottom)}]}),visibleItemsStyle:P(()=>(g.value,{transform:`translateY(${n(h.value.sum(v.value))})`})),viewportItems:b,listElRef:u,itemsElRef:D(null),scrollTo:x,handleListResize:M,handleListScroll:A,handleListWheel:j,handleItemResize:E}},render(){let{itemResizable:e,keyField:t,keyToIndex:n,visibleItemsTag:r}=this;return q(o,{onResize:this.handleListResize},{default:()=>{var i;return q(`div`,Y(this.$attrs,{class:[`v-vl`,this.showScrollbar&&`v-vl--show-scrollbar`],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:`listElRef`}),[this.items.length===0?(i=this.$slots).empty?.call(i):q(`div`,{ref:`itemsElRef`,class:`v-vl-items`,style:this.itemsStyle},[q(r,Object.assign({class:`v-vl-visible-items`,style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>{let{renderCol:r,renderItemWithCols:i}=this;return this.viewportItems.map(a=>{let s=a[t],c=n.get(s),l=r==null?void 0:q(Ae,{index:c,item:a}),u=i==null?void 0:q(Ae,{index:c,item:a}),d=this.$slots.default({item:a,renderedCols:l,renderedItemWithCols:u,index:c})[0];return e?q(o,{key:s,onResize:e=>this.handleItemResize(s,e)},{default:()=>d}):(d.key=s,d)})}})])])}})}});function Ne(t,n){n&&(y(()=>{let{value:r}=t;r&&e.registerHandler(r,n)}),_(t,(t,n)=>{n&&e.unregisterHandler(n)},{deep:!1}),w(()=>{let{value:n}=t;n&&e.unregisterHandler(n)}))}function Pe(e,t){if(!e)return;let n=document.createElement(`a`);n.href=e,t!==void 0&&(n.download=t),document.body.appendChild(n),n.click(),document.body.removeChild(n)}function Fe(e,t){Pe(e,t)}var Ie=new WeakSet;function Le(e){Ie.add(e)}function Re(e){return!Ie.has(e)}function ze(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var Be=R({name:`Checkmark`,render(){return q(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},q(`g`,{fill:`none`},q(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Ve=R({name:`Empty`,render(){return q(`svg`,{viewBox:`0 0 28 28`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},q(`path`,{d:`M26 7.5C26 11.0899 23.0899 14 19.5 14C15.9101 14 13 11.0899 13 7.5C13 3.91015 15.9101 1 19.5 1C23.0899 1 26 3.91015 26 7.5ZM16.8536 4.14645C16.6583 3.95118 16.3417 3.95118 16.1464 4.14645C15.9512 4.34171 15.9512 4.65829 16.1464 4.85355L18.7929 7.5L16.1464 10.1464C15.9512 10.3417 15.9512 10.6583 16.1464 10.8536C16.3417 11.0488 16.6583 11.0488 16.8536 10.8536L19.5 8.20711L22.1464 10.8536C22.3417 11.0488 22.6583 11.0488 22.8536 10.8536C23.0488 10.6583 23.0488 10.3417 22.8536 10.1464L20.2071 7.5L22.8536 4.85355C23.0488 4.65829 23.0488 4.34171 22.8536 4.14645C22.6583 3.95118 22.3417 3.95118 22.1464 4.14645L19.5 6.79289L16.8536 4.14645Z`,fill:`currentColor`}),q(`path`,{d:`M25 22.75V12.5991C24.5572 13.0765 24.053 13.4961 23.5 13.8454V16H17.5L17.3982 16.0068C17.0322 16.0565 16.75 16.3703 16.75 16.75C16.75 18.2688 15.5188 19.5 14 19.5C12.4812 19.5 11.25 18.2688 11.25 16.75L11.2432 16.6482C11.1935 16.2822 10.8797 16 10.5 16H4.5V7.25C4.5 6.2835 5.2835 5.5 6.25 5.5H12.2696C12.4146 4.97463 12.6153 4.47237 12.865 4H6.25C4.45507 4 3 5.45507 3 7.25V22.75C3 24.5449 4.45507 26 6.25 26H21.75C23.5449 26 25 24.5449 25 22.75ZM4.5 22.75V17.5H9.81597L9.85751 17.7041C10.2905 19.5919 11.9808 21 14 21L14.215 20.9947C16.2095 20.8953 17.842 19.4209 18.184 17.5H23.5V22.75C23.5 23.7165 22.7165 24.5 21.75 24.5H6.25C5.2835 24.5 4.5 23.7165 4.5 22.75Z`,fill:`currentColor`}))}}),He=R({props:{onFocus:Function,onBlur:Function},setup(e){return()=>q(`div`,{style:`width: 0; height: 0`,tabindex:0,onFocus:e.onFocus,onBlur:e.onBlur})}}),Ue={iconSizeTiny:`28px`,iconSizeSmall:`34px`,iconSizeMedium:`40px`,iconSizeLarge:`46px`,iconSizeHuge:`52px`};function We(e){let{textColorDisabled:t,iconColor:n,textColor2:r,fontSizeTiny:i,fontSizeSmall:a,fontSizeMedium:o,fontSizeLarge:s,fontSizeHuge:c}=e;return Object.assign(Object.assign({},Ue),{fontSizeTiny:i,fontSizeSmall:a,fontSizeMedium:o,fontSizeLarge:s,fontSizeHuge:c,textColor:t,iconColor:n,extraTextColor:r})}var Ge={name:`Empty`,common:V,self:We},Ke=h(`empty`,`
 display: flex;
 flex-direction: column;
 align-items: center;
 font-size: var(--n-font-size);
`,[E(`icon`,`
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 font-size: var(--n-icon-size);
 line-height: var(--n-icon-size);
 color: var(--n-icon-color);
 transition:
 color .3s var(--n-bezier);
 `,[T(`+`,[E(`description`,`
 margin-top: 8px;
 `)])]),E(`description`,`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 `),E(`extra`,`
 text-align: center;
 transition: color .3s var(--n-bezier);
 margin-top: 12px;
 color: var(--n-extra-text-color);
 `)]),qe=Object.assign(Object.assign({},O.props),{description:String,showDescription:{type:Boolean,default:!0},showIcon:{type:Boolean,default:!0},size:{type:String,default:`medium`},renderIcon:Function}),Je=R({name:`Empty`,props:qe,slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:r}=I(e),i=O(`Empty`,`-empty`,Ke,Ge,e,t),{localeRef:a}=be(`Empty`),o=P(()=>e.description??r?.value?.Empty?.description),s=P(()=>r?.value?.Empty?.renderIcon||(()=>q(Ve,null))),c=P(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{[U(`iconSize`,t)]:r,[U(`fontSize`,t)]:a,textColor:o,iconColor:s,extraTextColor:c}}=i.value;return{"--n-icon-size":r,"--n-font-size":a,"--n-bezier":n,"--n-text-color":o,"--n-icon-color":s,"--n-extra-text-color":c}}),l=n?u(`empty`,P(()=>{let t=``,{size:n}=e;return t+=n[0],t}),c,e):void 0;return{mergedClsPrefix:t,mergedRenderIcon:s,localizedDescription:P(()=>o.value||a.value.description),cssVars:n?void 0:c,themeClass:l?.themeClass,onRender:l?.onRender}},render(){let{$slots:e,mergedClsPrefix:t,onRender:n}=this;return n?.(),q(`div`,{class:[`${t}-empty`,this.themeClass],style:this.cssVars},this.showIcon?q(`div`,{class:`${t}-empty__icon`},e.icon?e.icon():q(z,{clsPrefix:t},{default:this.mergedRenderIcon})):null,this.showDescription?q(`div`,{class:`${t}-empty__description`},e.default?e.default():this.localizedDescription):null,e.extra?q(`div`,{class:`${t}-empty__extra`},e.extra()):null)}}),Ye={height:`calc(var(--n-option-height) * 7.6)`,paddingTiny:`4px 0`,paddingSmall:`4px 0`,paddingMedium:`4px 0`,paddingLarge:`4px 0`,paddingHuge:`4px 0`,optionPaddingTiny:`0 12px`,optionPaddingSmall:`0 12px`,optionPaddingMedium:`0 12px`,optionPaddingLarge:`0 12px`,optionPaddingHuge:`0 12px`,loadingSize:`18px`};function Xe(e){let{borderRadius:t,popoverColor:n,textColor3:r,dividerColor:i,textColor2:a,primaryColorPressed:o,textColorDisabled:s,primaryColor:c,opacityDisabled:l,hoverColor:u,fontSizeTiny:d,fontSizeSmall:f,fontSizeMedium:p,fontSizeLarge:m,fontSizeHuge:h,heightTiny:g,heightSmall:_,heightMedium:v,heightLarge:y,heightHuge:b}=e;return Object.assign(Object.assign({},Ye),{optionFontSizeTiny:d,optionFontSizeSmall:f,optionFontSizeMedium:p,optionFontSizeLarge:m,optionFontSizeHuge:h,optionHeightTiny:g,optionHeightSmall:_,optionHeightMedium:v,optionHeightLarge:y,optionHeightHuge:b,borderRadius:t,color:n,groupHeaderTextColor:r,actionDividerColor:i,optionTextColor:a,optionTextColorPressed:o,optionTextColorDisabled:s,optionTextColorActive:c,optionOpacityDisabled:l,optionCheckColor:c,optionColorPending:u,optionColorActive:`rgba(0, 0, 0, 0)`,optionColorActivePending:u,actionTextColor:a,loadingColor:c})}var Ze=L({name:`InternalSelectMenu`,common:V,peers:{Scrollbar:s,Empty:Ge},self:Xe}),Qe=R({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=M(le);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):ue(i[this.labelField],i,!1),s=q(`div`,Object.assign({},a,{class:[`${e}-base-select-group-header`,a?.class]}),o);return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function $e(e,t){return q(x,{name:`fade-in-scale-up-transition`},{default:()=>e?q(z,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>q(Be)}):null})}var et=R({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:f}=M(le),p=X(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function m(t){let{tmNode:n}=e;n.disabled||d(t,n)}function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e,{value:r}=p;n.disabled||r||f(t,n)}return{multiple:r,isGrouped:X(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:p,isSelected:X(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}else return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:g,handleMouseEnter:h,handleClick:m}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:a,nodeProps:o,renderOption:s,renderLabel:c,handleClick:l,handleMouseEnter:u,handleMouseMove:d}=this,f=$e(n,e),p=c?[c(t,n),a&&f]:[ue(t[this.labelField],t,n),a&&f],m=o?.(t),h=q(`div`,Object.assign({},m,{class:[`${e}-base-select-option`,t.class,m?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:a}],style:[m?.style||``,t.style||``],onClick:ze([l,m?.onClick]),onMouseenter:ze([u,m?.onMouseenter]),onMousemove:ze([d,m?.onMousemove])}),q(`div`,{class:`${e}-base-select-option__content`},p));return t.render?t.render({node:h,option:t,selected:n}):s?s({node:h,option:t,selected:n}):h}}),tt=h(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[h(`scrollbar`,`
 max-height: var(--n-height);
 `),h(`virtual-list`,`
 max-height: var(--n-height);
 `),h(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[E(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),h(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),h(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),E(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),E(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),E(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),E(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),h(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),h(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[m(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),T(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),T(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),m(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),m(`pending`,[T(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),m(`selected`,`
 color: var(--n-option-text-color-active);
 `,[T(`&::before`,`
 background-color: var(--n-option-color-active);
 `),m(`pending`,[T(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),m(`disabled`,`
 cursor: not-allowed;
 `,[g(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),m(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),E(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[me({enterScale:`0.5`})])])]),nt=R({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},O.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(e){let{mergedClsPrefixRef:n,mergedRtlRef:r,mergedComponentPropsRef:a}=I(e),o=j(`InternalSelectMenu`,r,n),s=O(`InternalSelectMenu`,`-internal-select-menu`,tt,Ze,e,F(e,`clsPrefix`)),c=D(null),l=D(null),d=D(null),f=P(()=>e.treeMate.getFlattenedNodes()),p=P(()=>he(f.value)),m=D(null);function h(){let{treeMate:t}=e,n=null,{value:r}=e;r===null?n=t.getFirstAvailableNode():(n=e.multiple?t.getNode((r||[])[(r||[]).length-1]):t.getNode(r),(!n||n.disabled)&&(n=t.getFirstAvailableNode())),G(n||null)}function g(){let{value:t}=m;t&&!e.treeMate.getNode(t.key)&&(m.value=null)}let b;_(()=>e.show,t=>{t?b=_(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?h():g(),C(K)):g()},{immediate:!0}):b?.()},{immediate:!0}),w(()=>{b?.()});let x=P(()=>i(s.value.self[U(`optionHeight`,e.size)])),S=P(()=>t(s.value.self[U(`padding`,e.size)])),T=P(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),E=P(()=>{let e=f.value;return e&&e.length===0}),k=P(()=>a?.value?.Select?.renderEmpty);function A(t){let{onToggle:n}=e;n&&n(t)}function M(t){let{onScroll:n}=e;n&&n(t)}function N(e){var t;(t=d.value)==null||t.sync(),M(e)}function L(){var e;(e=d.value)==null||e.sync()}function R(){let{value:e}=m;return e||null}function z(e,t){t.disabled||G(t,!1)}function ee(e,t){t.disabled||A(t)}function B(t){var n;ge(t,`action`)||(n=e.onKeyup)==null||n.call(e,t)}function V(t){var n;ge(t,`action`)||(n=e.onKeydown)==null||n.call(e,t)}function H(t){var n;(n=e.onMousedown)==null||n.call(e,t),!e.focusable&&t.preventDefault()}function te(){let{value:e}=m;e&&G(e.getNext({loop:!0}),!0)}function W(){let{value:e}=m;e&&G(e.getPrev({loop:!0}),!0)}function G(e,t=!1){m.value=e,t&&K()}function K(){var t,n;let r=m.value;if(!r)return;let i=p.value(r.key);i!==null&&(e.virtualScroll?(t=l.value)==null||t.scrollTo({index:i}):(n=d.value)==null||n.scrollTo({index:i,elSize:x.value}))}function q(t){var n;c.value?.contains(t.target)&&((n=e.onFocus)==null||n.call(e,t))}function J(t){var n;c.value?.contains(t.relatedTarget)||(n=e.onBlur)==null||n.call(e,t)}v(le,{handleOptionMouseEnter:z,handleOptionClick:ee,valueSetRef:T,pendingTmNodeRef:m,nodePropsRef:F(e,`nodeProps`),showCheckmarkRef:F(e,`showCheckmark`),multipleRef:F(e,`multiple`),valueRef:F(e,`value`),renderLabelRef:F(e,`renderLabel`),renderOptionRef:F(e,`renderOption`),labelFieldRef:F(e,`labelField`),valueFieldRef:F(e,`valueField`)}),v(re,c),y(()=>{let{value:e}=d;e&&e.sync()});let Y=P(()=>{let{size:n}=e,{common:{cubicBezierEaseInOut:r},self:{height:i,borderRadius:a,color:o,groupHeaderTextColor:c,actionDividerColor:l,optionTextColorPressed:u,optionTextColor:d,optionTextColorDisabled:f,optionTextColorActive:p,optionOpacityDisabled:m,optionCheckColor:h,actionTextColor:g,optionColorPending:_,optionColorActive:v,loadingColor:y,loadingSize:b,optionColorActivePending:x,[U(`optionFontSize`,n)]:S,[U(`optionHeight`,n)]:C,[U(`optionPadding`,n)]:w}}=s.value;return{"--n-height":i,"--n-action-divider-color":l,"--n-action-text-color":g,"--n-bezier":r,"--n-border-radius":a,"--n-color":o,"--n-option-font-size":S,"--n-group-header-text-color":c,"--n-option-check-color":h,"--n-option-color-pending":_,"--n-option-color-active":v,"--n-option-color-active-pending":x,"--n-option-height":C,"--n-option-opacity-disabled":m,"--n-option-text-color":d,"--n-option-text-color-active":p,"--n-option-text-color-disabled":f,"--n-option-text-color-pressed":u,"--n-option-padding":w,"--n-option-padding-left":t(w,`left`),"--n-option-padding-right":t(w,`right`),"--n-loading-color":y,"--n-loading-size":b}}),{inlineThemeDisabled:X}=e,ne=X?u(`internal-select-menu`,P(()=>e.size[0]),Y,e):void 0,Z={selfRef:c,next:te,prev:W,getPendingTmNode:R};return Ne(c,e.onResize),Object.assign({mergedTheme:s,mergedClsPrefix:n,rtlEnabled:o,virtualListRef:l,scrollbarRef:d,itemSize:x,padding:S,flattenedNodes:f,empty:E,mergedRenderEmpty:k,virtualListContainer(){let{value:e}=l;return e?.listElRef},virtualListContent(){let{value:e}=l;return e?.itemsElRef},doScroll:M,handleFocusin:q,handleFocusout:J,handleKeyUp:B,handleKeyDown:V,handleMouseDown:H,handleVirtualListResize:L,handleVirtualListScroll:N,cssVars:X?void 0:Y,themeClass:ne?.themeClass,onRender:ne?.onRender},Z)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:r,themeClass:i,onRender:a}=this;return a?.(),q(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,i,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},W(e.header,e=>e&&q(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?q(`div`,{class:`${n}-base-select-menu__loading`},q(K,{clsPrefix:n,strokeWidth:20})):this.empty?q(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},A(e.empty,()=>[this.mergedRenderEmpty?.call(this)||q(Je,{theme:r.peers.Empty,themeOverrides:r.peerOverrides.Empty,size:this.size})])):q(c,Object.assign({ref:`scrollbarRef`,theme:r.peers.Scrollbar,themeOverrides:r.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?q(Me,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?q(Qe,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:q(et,{clsPrefix:n,key:e.key,tmNode:e})}):q(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?q(Qe,{key:e.key,clsPrefix:n,tmNode:e}):q(et,{clsPrefix:n,key:e.key,tmNode:e})))}),W(e.action,e=>e&&[q(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),q(He,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),rt={closeIconSizeTiny:`12px`,closeIconSizeSmall:`12px`,closeIconSizeMedium:`14px`,closeIconSizeLarge:`14px`,closeSizeTiny:`16px`,closeSizeSmall:`16px`,closeSizeMedium:`18px`,closeSizeLarge:`18px`,padding:`0 7px`,closeMargin:`0 0 0 4px`};function it(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:c,baseColor:l,borderColor:u,opacityDisabled:d,tagColor:f,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,borderRadiusSmall:g,fontSizeMini:_,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,heightMini:x,heightTiny:S,heightSmall:C,heightMedium:w,closeColorHover:T,closeColorPressed:E,buttonColor2Hover:D,buttonColor2Pressed:O,fontWeightStrong:k}=e;return Object.assign(Object.assign({},rt),{closeBorderRadius:g,heightTiny:x,heightSmall:S,heightMedium:C,heightLarge:w,borderRadius:g,opacityDisabled:d,fontSizeTiny:_,fontSizeSmall:v,fontSizeMedium:y,fontSizeLarge:b,fontWeightStrong:k,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:l,colorCheckable:`#0000`,colorHoverCheckable:D,colorPressedCheckable:O,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${u}`,textColor:t,color:f,colorBordered:`rgb(250, 250, 252)`,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,closeColorHover:T,closeColorPressed:E,borderPrimary:`1px solid ${N(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:N(i,{alpha:.12}),colorBorderedPrimary:N(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:N(i,{alpha:.12}),closeColorPressedPrimary:N(i,{alpha:.18}),borderInfo:`1px solid ${N(a,{alpha:.3})}`,textColorInfo:a,colorInfo:N(a,{alpha:.12}),colorBorderedInfo:N(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:N(a,{alpha:.12}),closeColorPressedInfo:N(a,{alpha:.18}),borderSuccess:`1px solid ${N(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:N(o,{alpha:.12}),colorBorderedSuccess:N(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:N(o,{alpha:.12}),closeColorPressedSuccess:N(o,{alpha:.18}),borderWarning:`1px solid ${N(s,{alpha:.35})}`,textColorWarning:s,colorWarning:N(s,{alpha:.15}),colorBorderedWarning:N(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:N(s,{alpha:.12}),closeColorPressedWarning:N(s,{alpha:.18}),borderError:`1px solid ${N(c,{alpha:.23})}`,textColorError:c,colorError:N(c,{alpha:.1}),colorBorderedError:N(c,{alpha:.08}),closeIconColorError:c,closeIconColorHoverError:c,closeIconColorPressedError:c,closeColorHoverError:N(c,{alpha:.12}),closeColorPressedError:N(c,{alpha:.18})})}var at={name:`Tag`,common:V,self:it},ot={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},st=h(`tag`,`
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
`,[m(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),E(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),E(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),E(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),E(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),m(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[E(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),E(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),m(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),m(`icon, avatar`,[m(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),m(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),m(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[g(`disabled`,[T(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[g(`checked`,`color: var(--n-text-color-hover-checkable);`)]),T(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[g(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),m(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[g(`disabled`,[T(`&:hover`,`background-color: var(--n-color-checked-hover);`),T(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),ct=Object.assign(Object.assign(Object.assign({},O.props),ot),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),lt=k(`n-tag`),ut=R({name:`Tag`,props:ct,slots:Object,setup(e){let n=D(null),{mergedBorderedRef:r,mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedRtlRef:o,mergedComponentPropsRef:s}=I(e),c=P(()=>e.size||s?.value?.Tag?.size||`medium`),l=O(`Tag`,`-tag`,st,at,e,i);v(lt,{roundRef:F(e,`round`)});function d(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function f(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&H(n,t)}}let p={setTextContent(e){let{value:t}=n;t&&(t.textContent=e)}},m=j(`Tag`,o,i),h=P(()=>{let{type:n,color:{color:i,textColor:a}={}}=e,o=c.value,{common:{cubicBezierEaseInOut:s},self:{padding:u,closeMargin:d,borderRadius:f,opacityDisabled:p,textColorCheckable:m,textColorHoverCheckable:h,textColorPressedCheckable:g,textColorChecked:_,colorCheckable:v,colorHoverCheckable:y,colorPressedCheckable:b,colorChecked:x,colorCheckedHover:S,colorCheckedPressed:C,closeBorderRadius:w,fontWeightStrong:T,[U(`colorBordered`,n)]:E,[U(`closeSize`,o)]:D,[U(`closeIconSize`,o)]:O,[U(`fontSize`,o)]:k,[U(`height`,o)]:A,[U(`color`,n)]:j,[U(`textColor`,n)]:M,[U(`border`,n)]:N,[U(`closeIconColor`,n)]:P,[U(`closeIconColorHover`,n)]:F,[U(`closeIconColorPressed`,n)]:I,[U(`closeColorHover`,n)]:L,[U(`closeColorPressed`,n)]:R}}=l.value,z=t(d);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${A} - 8px)`,"--n-bezier":s,"--n-border-radius":f,"--n-border":N,"--n-close-icon-size":O,"--n-close-color-pressed":R,"--n-close-color-hover":L,"--n-close-border-radius":w,"--n-close-icon-color":P,"--n-close-icon-color-hover":F,"--n-close-icon-color-pressed":I,"--n-close-icon-color-disabled":P,"--n-close-margin-top":z.top,"--n-close-margin-right":z.right,"--n-close-margin-bottom":z.bottom,"--n-close-margin-left":z.left,"--n-close-size":D,"--n-color":i||(r.value?E:j),"--n-color-checkable":v,"--n-color-checked":x,"--n-color-checked-hover":S,"--n-color-checked-pressed":C,"--n-color-hover-checkable":y,"--n-color-pressed-checkable":b,"--n-font-size":k,"--n-height":A,"--n-opacity-disabled":p,"--n-padding":u,"--n-text-color":a||M,"--n-text-color-checkable":m,"--n-text-color-checked":_,"--n-text-color-hover-checkable":h,"--n-text-color-pressed-checkable":g}}),g=a?u(`tag`,P(()=>{let t=``,{type:n,color:{color:i,textColor:a}={}}=e;return t+=n[0],t+=c.value[0],i&&(t+=`a${te(i)}`),a&&(t+=`b${te(a)}`),r.value&&(t+=`c`),t}),h,e):void 0;return Object.assign(Object.assign({},p),{rtlEnabled:m,mergedClsPrefix:i,contentRef:n,mergedBordered:r,handleClick:d,handleCloseClick:f,cssVars:a?void 0:h,themeClass:g?.themeClass,onRender:g?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:a,onRender:o,$slots:s}=this;o?.();let c=W(s.avatar,e=>e&&q(`div`,{class:`${t}-tag__avatar`},e)),l=W(s.icon,e=>e&&q(`div`,{class:`${t}-tag__icon`},e));return q(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:a,[`${t}-tag--avatar`]:c,[`${t}-tag--icon`]:l,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},l||c,q(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?q(xe,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:a,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?q(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),dt={paddingSingle:`0 26px 0 12px`,paddingMultiple:`3px 26px 0 12px`,clearSize:`16px`,arrowSize:`16px`};function ft(e){let{borderRadius:t,textColor2:n,textColorDisabled:r,inputColor:i,inputColorDisabled:a,primaryColor:o,primaryColorHover:s,warningColor:c,warningColorHover:l,errorColor:u,errorColorHover:d,borderColor:f,iconColor:p,iconColorDisabled:m,clearColor:h,clearColorHover:g,clearColorPressed:_,placeholderColor:v,placeholderColorDisabled:y,fontSizeTiny:b,fontSizeSmall:x,fontSizeMedium:S,fontSizeLarge:C,heightTiny:w,heightSmall:T,heightMedium:E,heightLarge:D,fontWeight:O}=e;return Object.assign(Object.assign({},dt),{fontSizeTiny:b,fontSizeSmall:x,fontSizeMedium:S,fontSizeLarge:C,heightTiny:w,heightSmall:T,heightMedium:E,heightLarge:D,borderRadius:t,fontWeight:O,textColor:n,textColorDisabled:r,placeholderColor:v,placeholderColorDisabled:y,color:i,colorDisabled:a,colorActive:i,border:`1px solid ${f}`,borderHover:`1px solid ${s}`,borderActive:`1px solid ${o}`,borderFocus:`1px solid ${s}`,boxShadowHover:`none`,boxShadowActive:`0 0 0 2px ${N(o,{alpha:.2})}`,boxShadowFocus:`0 0 0 2px ${N(o,{alpha:.2})}`,caretColor:o,arrowColor:p,arrowColorDisabled:m,loadingColor:o,borderWarning:`1px solid ${c}`,borderHoverWarning:`1px solid ${l}`,borderActiveWarning:`1px solid ${c}`,borderFocusWarning:`1px solid ${l}`,boxShadowHoverWarning:`none`,boxShadowActiveWarning:`0 0 0 2px ${N(c,{alpha:.2})}`,boxShadowFocusWarning:`0 0 0 2px ${N(c,{alpha:.2})}`,colorActiveWarning:i,caretColorWarning:c,borderError:`1px solid ${u}`,borderHoverError:`1px solid ${d}`,borderActiveError:`1px solid ${u}`,borderFocusError:`1px solid ${d}`,boxShadowHoverError:`none`,boxShadowActiveError:`0 0 0 2px ${N(u,{alpha:.2})}`,boxShadowFocusError:`0 0 0 2px ${N(u,{alpha:.2})}`,colorActiveError:i,caretColorError:u,clearColor:h,clearColorHover:g,clearColorPressed:_})}var pt=L({name:`InternalSelection`,common:V,peers:{Popover:de},self:ft}),mt=T([h(`base-selection`,`
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
 `,[h(`base-loading`,`
 color: var(--n-loading-color);
 `),h(`base-selection-tags`,`min-height: var(--n-height);`),E(`border, state-border`,`
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
 `),E(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),h(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[E(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),h(`base-selection-overlay`,`
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
 `,[E(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),h(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[E(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),h(`base-selection-tags`,`
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
 `),h(`base-selection-label`,`
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
 `,[h(`base-selection-input`,`
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
 `,[E(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),E(`render-label`,`
 color: var(--n-text-color);
 `)]),g(`disabled`,[T(`&:hover`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),m(`focus`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),m(`active`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),h(`base-selection-label`,`background-color: var(--n-color-active);`),h(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),m(`disabled`,`cursor: not-allowed;`,[E(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),h(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[h(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),E(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),h(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),h(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),h(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[E(`input`,`
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
 `),E(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>m(`${e}-status`,[E(`state-border`,`border: var(--n-border-${e});`),g(`disabled`,[T(`&:hover`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),m(`active`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),h(`base-selection-label`,`background-color: var(--n-color-active-${e});`),h(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),m(`focus`,[E(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),h(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),h(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[T(`&:last-child`,`padding-right: 0;`),h(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[E(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),ht=R({name:`InternalSelection`,props:Object.assign(Object.assign({},O.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(e){let{mergedClsPrefixRef:n,mergedRtlRef:r}=I(e),i=j(`InternalSelection`,r,n),a=D(null),o=D(null),s=D(null),c=D(null),l=D(null),d=D(null),f=D(null),p=D(null),m=D(null),h=D(null),g=D(!1),v=D(!1),b=D(!1),x=O(`InternalSelection`,`-internal-selection`,mt,pt,e,F(e,`clsPrefix`)),w=P(()=>e.clearable&&!e.disabled&&(b.value||e.active)),T=P(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):ue(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),E=P(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),k=P(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function A(){var t;let{value:n}=a;if(n){let{value:r}=o;r&&(r.style.width=`${n.offsetWidth}px`,e.maxTagCount!==`responsive`&&((t=m.value)==null||t.sync({showAllItemsBeforeCalculate:!1})))}}function M(){let{value:e}=h;e&&(e.style.display=`none`)}function N(){let{value:e}=h;e&&(e.style.display=`inline-block`)}_(F(e,`active`),e=>{e||M()}),_(F(e,`pattern`),()=>{e.multiple&&C(A)});function L(t){let{onFocus:n}=e;n&&n(t)}function R(t){let{onBlur:n}=e;n&&n(t)}function z(t){let{onDeleteOption:n}=e;n&&n(t)}function ee(t){let{onClear:n}=e;n&&n(t)}function B(t){let{onPatternInput:n}=e;n&&n(t)}function V(e){(!e.relatedTarget||!s.value?.contains(e.relatedTarget))&&L(e)}function H(e){s.value?.contains(e.relatedTarget)||R(e)}function te(e){ee(e)}function W(){b.value=!0}function G(){b.value=!1}function K(t){!e.active||!e.filterable||t.target!==o.value&&t.preventDefault()}function q(e){z(e)}let J=D(!1);function Y(t){if(t.key===`Backspace`&&!J.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&q(t[t.length-1])}}let X=null;function ne(t){let{value:n}=a;n&&(n.textContent=t.target.value,A()),e.ignoreComposition&&J.value?X=t:B(t)}function re(){J.value=!0}function Z(){J.value=!1,e.ignoreComposition&&B(X),X=null}function ie(t){var n;v.value=!0,(n=e.onPatternFocus)==null||n.call(e,t)}function Q(t){var n;v.value=!1,(n=e.onPatternBlur)==null||n.call(e,t)}function $(){var t,n;if(e.filterable)v.value=!1,(t=d.value)==null||t.blur(),(n=o.value)==null||n.blur();else if(e.multiple){let{value:e}=c;e?.blur()}else{let{value:e}=l;e?.blur()}}function ae(){var t,n,r;e.filterable?(v.value=!1,(t=d.value)==null||t.focus()):e.multiple?(n=c.value)==null||n.focus():(r=l.value)==null||r.focus()}function oe(){let{value:e}=o;e&&(N(),e.focus())}function se(){let{value:e}=o;e&&e.blur()}function ce(e){let{value:t}=f;t&&t.setTextContent(`+${e}`)}function le(){let{value:e}=p;return e}function de(){return o.value}let fe=null;function pe(){fe!==null&&window.clearTimeout(fe)}function me(){e.active||(pe(),fe=window.setTimeout(()=>{k.value&&(g.value=!0)},100))}function he(){pe()}function ge(e){e||(pe(),g.value=!1)}_(k,e=>{e||(g.value=!1)}),y(()=>{S(()=>{let t=d.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=v.value?-1:0)})}),Ne(s,e.onResize);let{inlineThemeDisabled:_e}=e,ve=P(()=>{let{size:n}=e,{common:{cubicBezierEaseInOut:r},self:{fontWeight:i,borderRadius:a,color:o,placeholderColor:s,textColor:c,paddingSingle:l,paddingMultiple:u,caretColor:d,colorDisabled:f,textColorDisabled:p,placeholderColorDisabled:m,colorActive:h,boxShadowFocus:g,boxShadowActive:_,boxShadowHover:v,border:y,borderFocus:b,borderHover:S,borderActive:C,arrowColor:w,arrowColorDisabled:T,loadingColor:E,colorActiveWarning:D,boxShadowFocusWarning:O,boxShadowActiveWarning:k,boxShadowHoverWarning:A,borderWarning:j,borderFocusWarning:M,borderHoverWarning:N,borderActiveWarning:P,colorActiveError:F,boxShadowFocusError:I,boxShadowActiveError:L,boxShadowHoverError:R,borderError:z,borderFocusError:ee,borderHoverError:B,borderActiveError:V,clearColor:H,clearColorHover:te,clearColorPressed:W,clearSize:G,arrowSize:K,[U(`height`,n)]:q,[U(`fontSize`,n)]:J}}=x.value,Y=t(l),X=t(u);return{"--n-bezier":r,"--n-border":y,"--n-border-active":C,"--n-border-focus":b,"--n-border-hover":S,"--n-border-radius":a,"--n-box-shadow-active":_,"--n-box-shadow-focus":g,"--n-box-shadow-hover":v,"--n-caret-color":d,"--n-color":o,"--n-color-active":h,"--n-color-disabled":f,"--n-font-size":J,"--n-height":q,"--n-padding-single-top":Y.top,"--n-padding-multiple-top":X.top,"--n-padding-single-right":Y.right,"--n-padding-multiple-right":X.right,"--n-padding-single-left":Y.left,"--n-padding-multiple-left":X.left,"--n-padding-single-bottom":Y.bottom,"--n-padding-multiple-bottom":X.bottom,"--n-placeholder-color":s,"--n-placeholder-color-disabled":m,"--n-text-color":c,"--n-text-color-disabled":p,"--n-arrow-color":w,"--n-arrow-color-disabled":T,"--n-loading-color":E,"--n-color-active-warning":D,"--n-box-shadow-focus-warning":O,"--n-box-shadow-active-warning":k,"--n-box-shadow-hover-warning":A,"--n-border-warning":j,"--n-border-focus-warning":M,"--n-border-hover-warning":N,"--n-border-active-warning":P,"--n-color-active-error":F,"--n-box-shadow-focus-error":I,"--n-box-shadow-active-error":L,"--n-box-shadow-hover-error":R,"--n-border-error":z,"--n-border-focus-error":ee,"--n-border-hover-error":B,"--n-border-active-error":V,"--n-clear-size":G,"--n-clear-color":H,"--n-clear-color-hover":te,"--n-clear-color-pressed":W,"--n-arrow-size":K,"--n-font-weight":i}}),ye=_e?u(`internal-selection`,P(()=>e.size[0]),ve,e):void 0;return{mergedTheme:x,mergedClearable:w,mergedClsPrefix:n,rtlEnabled:i,patternInputFocused:v,filterablePlaceholder:T,label:E,selected:k,showTagsPanel:g,isComposing:J,counterRef:f,counterWrapperRef:p,patternInputMirrorRef:a,patternInputRef:o,selfRef:s,multipleElRef:c,singleElRef:l,patternInputWrapperRef:d,overflowRef:m,inputTagElRef:h,handleMouseDown:K,handleFocusin:V,handleClear:te,handleMouseEnter:W,handleMouseLeave:G,handleDeleteOption:q,handlePatternKeyDown:Y,handlePatternInputInput:ne,handlePatternInputBlur:Q,handlePatternInputFocus:ie,handleMouseEnterCounter:me,handleMouseLeaveCounter:he,handleFocusout:H,handleCompositionEnd:Z,handleCompositionStart:re,onPopoverUpdateShow:ge,focus:ae,focusInput:oe,blur:$,blurInput:se,updateCounter:ce,getCounter:le,getTail:de,renderLabel:e.renderLabel,cssVars:_e?void 0:ve,themeClass:ye?.themeClass,onRender:ye?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:i,maxTagCount:o,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:u,renderTag:d,renderLabel:f}=this;u?.();let p=o===`responsive`,m=typeof o==`number`,h=p||m,g=q(a,null,{default:()=>q(ye,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),_;if(t){let{labelField:e}=this,t=t=>q(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},d?d({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):q(ut,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>f?f(t,!0):ue(t[e],t,!0)})),a=()=>(m?this.selectedOptions.slice(0,o):this.selectedOptions).map(t),s=i?q(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},q(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),q(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,u=p?()=>q(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},q(ut,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,v;if(m){let e=this.selectedOptions.length-o;e>0&&(v=q(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},q(ut,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let y=p?i?q(ve,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:a,counter:u,tail:()=>s}):q(ve,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:a,counter:u}):m&&v?a().concat(v):a(),b=h?()=>q(`div`,{class:`${c}-base-selection-popover`},p?a():this.selectedOptions.map(t)):void 0,x=h?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,S=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?q(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},q(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,C=i?q(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},y,p?null:s,g):q(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},y,g);_=q(B,null,h?q(fe,Object.assign({},x,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>C,default:b}):C,S)}else if(i){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;_=q(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:se(this.label)},q(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?q(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},q(`div`,{class:`${c}-base-selection-overlay__wrapper`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):ue(this.label,this.selectedOption,!0))):null,t?q(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},q(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,g)}else _=q(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?q(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},q(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):q(`div`,{class:`${c}-base-selection-input`,title:se(this.label),key:`input`},q(`div`,{class:`${c}-base-selection-input__content`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):ue(this.label,this.selectedOption,!0))),g);return q(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},_,s?q(`div`,{class:`${c}-base-selection__border`}):null,s?q(`div`,{class:`${c}-base-selection__state-border`}):null)}});function gt(e){return e.type===`group`}function _t(e){return e.type===`ignored`}function vt(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function yt(e,t){return{getIsGroup:gt,getIgnored:_t,getKey(t){return gt(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function bt(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(gt(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(_t(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function xt(e,t,n){let r=new Map;return e.forEach(e=>{gt(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}function St(e){let{boxShadow2:t}=e;return{menuBoxShadow:t}}var Ct=L({name:`Select`,common:V,peers:{InternalSelection:pt,InternalSelectMenu:Ze},self:St}),wt=T([h(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),h(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[me({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Tt=Object.assign(Object.assign({},O.props),{to:ae.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),Et=R({name:`Select`,props:Tt,slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:o}=I(e),s=O(`Select`,`-select`,wt,Ct,e,t),c=D(e.defaultValue),d=r(F(e,`value`),c),f=D(!1),p=D(``),m=ie(e,[`items`,`options`]),h=D([]),g=D([]),v=P(()=>g.value.concat(h.value).concat(m.value)),y=P(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return vt(e,i);let a=t[r];return typeof a==`string`?vt(e,a):typeof a==`number`?vt(e,String(a)):!1}}),x=P(()=>{if(e.remote)return m.value;{let{value:t}=v,{value:n}=p;return!n.length||!e.filterable?t:bt(t,y.value,n,e.childrenField)}}),S=P(()=>{let{valueField:t,childrenField:n}=e,r=yt(t,n);return _e(x.value,r)}),C=P(()=>xt(v.value,e.valueField,e.childrenField)),w=D(!1),T=r(F(e,`show`),w),E=D(null),k=D(null),A=D(null),{localeRef:j}=be(`Select`),M=P(()=>e.placeholder??j.value.placeholder),N=[],L=D(new Map),R=P(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function z(t){let n=e.remote,{value:r}=L,{value:i}=C,{value:a}=R,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let ee=P(()=>{if(e.multiple){let{value:e}=d;return Array.isArray(e)?z(e):[]}return null}),B=P(()=>{let{value:t}=d;return!e.multiple&&!Array.isArray(t)?t===null?null:z([t])[0]||null:null}),V=b(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:o?.value?.Select?.size||`medium`}}),{mergedSizeRef:te,mergedDisabledRef:U,mergedStatusRef:W}=V;function G(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=V;r&&H(r,t,n),a&&H(a,t,n),i&&H(i,t,n),c.value=t,o(),s()}function K(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=V;n&&H(n,t),r()}function q(){let{onClear:t}=e;t&&H(t)}function Y(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=V;n&&H(n,t),i(),r&&Q()}function X(t){let{onSearch:n}=e;n&&H(n,t)}function ne(t){let{onScroll:n}=e;n&&H(n,t)}function re(){var t;let{remote:n,multiple:r}=e;if(n){let{value:n}=L;if(r){let{valueField:r}=e;(t=ee.value)==null||t.forEach(e=>{n.set(e[r],e)})}else{let t=B.value;t&&n.set(t[e.valueField],t)}}}function Z(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&H(n,t),r&&H(r,t),w.value=t}function Q(){U.value||(Z(!0),w.value=!0,e.filterable&&Ae())}function $(){Z(!1)}function oe(){p.value=``,g.value=N}let se=D(!1);function ce(){e.filterable&&(se.value=!0)}function le(){e.filterable&&(se.value=!1,T.value||oe())}function ue(){U.value||(T.value?e.filterable?Ae():$():Q())}function de(e){(A.value?.selfRef)?.contains(e.relatedTarget)||(f.value=!1,K(e),$())}function fe(e){Y(e),f.value=!0}function pe(){f.value=!0}function me(e){E.value?.$el.contains(e.relatedTarget)||(f.value=!1,K(e),$())}function he(){var e;(e=E.value)==null||e.focus(),$()}function ve(e){T.value&&(E.value?.$el.contains(l(e))||$())}function ye(t){if(!Array.isArray(t))return[];if(R.value)return Array.from(t);{let{remote:n}=e,{value:r}=C;if(n){let{value:e}=L;return t.filter(t=>r.has(t)||e.has(t))}else return t.filter(e=>r.has(e))}}function xe(e){Se(e.rawNode)}function Se(t){if(U.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=g,t=e[0]||null;if(t){let e=h.value;e.length?e.push(t):h.value=[t],g.value=N}}if(r&&L.value.set(t[a],t),e.multiple){let e=ye(d.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Ce(t[a]);~e&&(h.value.splice(e,1),i&&(p.value=``))}}else e.push(t[a]),i&&(p.value=``);G(e,z(e))}else{if(n&&!r){let e=Ce(t[a]);~e?h.value=[h.value[e]]:h.value=N}ke(),$(),G(t[a],t)}}function Ce(t){return h.value.findIndex(n=>n[e.valueField]===t)}function we(t){T.value||Q();let{value:n}=t.target;p.value=n;let{tag:r,remote:i}=e;if(X(n),r&&!i){if(!n){g.value=N;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;m.value.some(e=>e[i]===r[i]||e[a]===r[a])||h.value.some(e=>e[i]===r[i]||e[a]===r[a])?g.value=N:g.value=[r]}}function Te(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&$(),r&&!i&&a&&(h.value=N),q(),n?G([],[]):G(null,null)}function Ee(e){!ge(e,`action`)&&!ge(e,`empty`)&&!ge(e,`header`)&&e.preventDefault()}function De(e){ne(e)}function Oe(t){var n,r,i;if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!E.value?.isComposing){if(T.value){let t=A.value?.getPendingTmNode();t?xe(t):e.filterable||($(),ke())}else if(Q(),e.tag&&se.value){let t=g.value[0];if(t){let n=t[e.valueField],{value:r}=d;e.multiple&&Array.isArray(r)&&r.includes(n)||Se(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;T.value&&((n=A.value)==null||n.prev());break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;T.value?(r=A.value)==null||r.next():Q();break;case`Escape`:T.value&&(Le(t),$()),(i=E.value)==null||i.focus();break}}function ke(){var e;(e=E.value)==null||e.focus()}function Ae(){var e;(e=E.value)==null||e.focusInput()}function je(){var e;T.value&&((e=k.value)==null||e.syncPosition())}re(),_(F(e,`options`),re);let Me={focus:()=>{var e;(e=E.value)==null||e.focus()},focusInput:()=>{var e;(e=E.value)==null||e.focusInput()},blur:()=>{var e;(e=E.value)==null||e.blur()},blurInput:()=>{var e;(e=E.value)==null||e.blurInput()}},Ne=P(()=>{let{self:{menuBoxShadow:e}}=s.value;return{"--n-menu-box-shadow":e}}),Pe=a?u(`select`,void 0,Ne,e):void 0;return Object.assign(Object.assign({},Me),{mergedStatus:W,mergedClsPrefix:t,mergedBordered:n,namespace:i,treeMate:S,isMounted:J(),triggerRef:E,menuRef:A,pattern:p,uncontrolledShow:w,mergedShow:T,adjustedTo:ae(e),uncontrolledValue:c,mergedValue:d,followerRef:k,localizedPlaceholder:M,selectedOption:B,selectedOptions:ee,mergedSize:te,mergedDisabled:U,focused:f,activeWithoutMenuOpen:se,inlineThemeDisabled:a,onTriggerInputFocus:ce,onTriggerInputBlur:le,handleTriggerOrMenuResize:je,handleMenuFocus:pe,handleMenuBlur:me,handleMenuTabOut:he,handleTriggerClick:ue,handleToggle:xe,handleDeleteOption:Se,handlePatternInput:we,handleClear:Te,handleTriggerBlur:de,handleTriggerFocus:fe,handleKeydown:Oe,handleMenuAfterLeave:oe,handleMenuClickOutside:ve,handleMenuScroll:De,handleMenuKeydown:Oe,handleMenuMousedown:Ee,mergedTheme:s,cssVars:a?void 0:Ne,themeClass:Pe?.themeClass,onRender:Pe?.onRender})},render(){return q(`div`,{class:`${this.mergedClsPrefix}-select`},q($,null,{default:()=>[q(pe,null,{default:()=>q(ht,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),q(ce,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===ae.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>q(x,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),d(q(nt,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[G,this.mergedShow],[ne,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[ne,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{Fe as C,Pe as S,Me as T,He as _,ht as a,Re as b,lt as c,at as d,nt as f,Ge as g,qe as h,yt as i,ct as l,Je as m,Tt as n,pt as o,Ze as p,Ct as r,ut as s,Et as t,ot as u,Be as v,Ne as w,Le as x,ze as y};