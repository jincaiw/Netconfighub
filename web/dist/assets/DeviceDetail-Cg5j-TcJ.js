import{m as e,o as t}from"./api-CZuskCjY.js";import{$t as n,An as r,Bt as i,Ct as a,En as o,Et as s,Fn as c,Gt as l,Ht as u,J as d,Jt as f,Lt as p,Mn as m,Mt as h,Nn as g,Pt as _,Q as v,Tt as y,Ut as b,V as x,Vt as S,W as C,Wt as w,Yt as T,Z as E,Zt as D,_n as O,at as k,bt as A,cn as j,ct as M,dt as N,en as P,g as F,gn as I,h as L,hn as ee,ln as te,m as R,mn as z,on as ne,ot as B,pt as V,qt as H,s as U,t as W,tn as G,un as re,wt as K,xt as q,y as J}from"./Button-v7qKim5p.js";import{b as ie,m as ae,p as Y,x as oe}from"./Tooltip-qMw84P-o.js";import{E as se,_ as X,m as ce,w as le}from"./token-dQuugHoF.js";import{i as ue,s as de,t as Z}from"./Select-Cc8dF3j7.js";import{t as fe}from"./DataTable-tfB4mMMi.js";import{n as pe,t as Q}from"./DescriptionsItem-DX6xAaYV.js";import{t as me}from"./Add-DnwwuZfa.js";import{t as he}from"./Space-BzYw_YS1.js";import{t as ge}from"./Spin-lY4kzDV6.js";import{t as _e}from"./text-D6Dk2Q8D.js";import{a as $,f as ve,g as ye,n as be,p as xe,r as Se,t as Ce}from"./index-CESIsNoN.js";import{t as we}from"./_plugin-vue_export-helper-BDNMzG2s.js";import{i as Te,o as Ee,r as De}from"./device-CFc-hf6-.js";import{i as Oe,r as ke}from"./backup-BVVA-03O.js";import{a as Ae,i as je,r as Me,t as Ne}from"./config-AzIoV5pK.js";var Pe=Y(`.v-x-scroll`,{overflow:`auto`,scrollbarWidth:`none`},[Y(`&::-webkit-scrollbar`,{width:0,height:0})]),Fe=T({name:`XScroll`,props:{disabled:Boolean,onScroll:Function},setup(){let e=o(null);function t(e){!(e.currentTarget.offsetWidth<e.currentTarget.scrollWidth)||e.deltaY===0||(e.currentTarget.scrollLeft+=e.deltaY+e.deltaX,e.preventDefault())}let n=N();return Pe.mount({id:`vueuc/x-scroll`,head:!0,anchorMetaName:ae,ssr:n}),Object.assign({selfRef:e,handleWheel:t},{scrollTo(...t){var n;(n=e.value)==null||n.scrollTo(...t)}})},render(){return D(`div`,{ref:`selfRef`,onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:`v-x-scroll`},this.$slots)}}),Ie=/\s/;function Le(e){for(var t=e.length;t--&&Ie.test(e.charAt(t)););return t}var Re=/^\s+/;function ze(e){return e&&e.slice(0,Le(e)+1).replace(Re,``)}var Be=NaN,Ve=/^[-+]0x[0-9a-f]+$/i,He=/^0b[01]+$/i,Ue=/^0o[0-7]+$/i,We=parseInt;function Ge(e){if(typeof e==`number`)return e;if(C(e))return Be;if(x(e)){var t=typeof e.valueOf==`function`?e.valueOf():e;e=x(t)?t+``:t}if(typeof e!=`string`)return e===0?e:+e;e=ze(e);var n=He.test(e);return n||Ue.test(e)?We(e.slice(2),n?2:8):Ve.test(e)?Be:+e}var Ke=function(){return d.Date.now()},qe=`Expected a function`,Je=Math.max,Ye=Math.min;function Xe(e,t,n){var r,i,a,o,s,c,l=0,u=!1,d=!1,f=!0;if(typeof e!=`function`)throw TypeError(qe);t=Ge(t)||0,x(n)&&(u=!!n.leading,d=`maxWait`in n,a=d?Je(Ge(n.maxWait)||0,t):a,f=`trailing`in n?!!n.trailing:f);function p(t){var n=r,a=i;return r=i=void 0,l=t,o=e.apply(a,n),o}function m(e){return l=e,s=setTimeout(_,t),u?p(e):o}function h(e){var n=e-c,r=e-l,i=t-n;return d?Ye(i,a-r):i}function g(e){var n=e-c,r=e-l;return c===void 0||n>=t||n<0||d&&r>=a}function _(){var e=Ke();if(g(e))return v(e);s=setTimeout(_,h(e))}function v(e){return s=void 0,f&&r?p(e):(r=i=void 0,o)}function y(){s!==void 0&&clearTimeout(s),l=0,r=c=i=s=void 0}function b(){return s===void 0?o:v(Ke())}function S(){var e=Ke(),n=g(e);if(r=arguments,i=this,c=e,n){if(s===void 0)return m(c);if(d)return clearTimeout(s),s=setTimeout(_,t),p(c)}return s===void 0&&(s=setTimeout(_,t)),o}return S.cancel=y,S.flush=b,S}var Ze=`Expected a function`;function Qe(e,t,n){var r=!0,i=!0;if(typeof e!=`function`)throw TypeError(Ze);return x(n)&&(r=`leading`in n?!!n.leading:r,i=`trailing`in n?!!n.trailing:i),Xe(e,t,{leading:r,maxWait:t,trailing:i})}var $e=T({name:`ArrowBack`,render(){return D(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 24 24`},D(`path`,{d:`M0 0h24v24H0V0z`,fill:`none`}),D(`path`,{d:`M19 11H7.83l4.88-4.88c.39-.39.39-1.03 0-1.42-.39-.39-1.02-.39-1.41 0l-6.59 6.59c-.39.39-.39 1.02 0 1.41l6.59 6.59c.39.39 1.02.39 1.41 0 .39-.39.39-1.02 0-1.41L7.83 13H19c.55 0 1-.45 1-1s-.45-1-1-1z`}))}}),et={titleFontSize:`18px`,backSize:`22px`};function tt(e){let{textColor1:t,textColor2:n,textColor3:r,fontSize:i,fontWeightStrong:a,primaryColorHover:o,primaryColorPressed:s}=e;return Object.assign(Object.assign({},et),{titleFontWeight:a,fontSize:i,titleTextColor:t,backColor:n,backColorHover:o,backColorPressed:s,subtitleTextColor:r})}var nt=L({name:`PageHeader`,common:U,self:tt}),rt={tabFontSizeSmall:`14px`,tabFontSizeMedium:`14px`,tabFontSizeLarge:`16px`,tabGapSmallLine:`36px`,tabGapMediumLine:`36px`,tabGapLargeLine:`36px`,tabGapSmallLineVertical:`8px`,tabGapMediumLineVertical:`8px`,tabGapLargeLineVertical:`8px`,tabPaddingSmallLine:`6px 0`,tabPaddingMediumLine:`10px 0`,tabPaddingLargeLine:`14px 0`,tabPaddingVerticalSmallLine:`6px 12px`,tabPaddingVerticalMediumLine:`8px 16px`,tabPaddingVerticalLargeLine:`10px 20px`,tabGapSmallBar:`36px`,tabGapMediumBar:`36px`,tabGapLargeBar:`36px`,tabGapSmallBarVertical:`8px`,tabGapMediumBarVertical:`8px`,tabGapLargeBarVertical:`8px`,tabPaddingSmallBar:`4px 0`,tabPaddingMediumBar:`6px 0`,tabPaddingLargeBar:`10px 0`,tabPaddingVerticalSmallBar:`6px 12px`,tabPaddingVerticalMediumBar:`8px 16px`,tabPaddingVerticalLargeBar:`10px 20px`,tabGapSmallCard:`4px`,tabGapMediumCard:`4px`,tabGapLargeCard:`4px`,tabGapSmallCardVertical:`4px`,tabGapMediumCardVertical:`4px`,tabGapLargeCardVertical:`4px`,tabPaddingSmallCard:`8px 16px`,tabPaddingMediumCard:`10px 20px`,tabPaddingLargeCard:`12px 24px`,tabPaddingSmallSegment:`4px 0`,tabPaddingMediumSegment:`6px 0`,tabPaddingLargeSegment:`8px 0`,tabPaddingVerticalLargeSegment:`0 8px`,tabPaddingVerticalSmallCard:`8px 12px`,tabPaddingVerticalMediumCard:`10px 16px`,tabPaddingVerticalLargeCard:`12px 20px`,tabPaddingVerticalSmallSegment:`0 4px`,tabPaddingVerticalMediumSegment:`0 6px`,tabGapSmallSegment:`0`,tabGapMediumSegment:`0`,tabGapLargeSegment:`0`,tabGapSmallSegmentVertical:`0`,tabGapMediumSegmentVertical:`0`,tabGapLargeSegmentVertical:`0`,panePaddingSmall:`8px 0 0 0`,panePaddingMedium:`12px 0 0 0`,panePaddingLarge:`16px 0 0 0`,closeSize:`18px`,closeIconSize:`14px`};function it(e){let{textColor2:t,primaryColor:n,textColorDisabled:r,closeIconColor:i,closeIconColorHover:a,closeIconColorPressed:o,closeColorHover:s,closeColorPressed:c,tabColor:l,baseColor:u,dividerColor:d,fontWeight:f,textColor1:p,borderRadius:m,fontSize:h,fontWeightStrong:g}=e;return Object.assign(Object.assign({},rt),{colorSegment:l,tabFontSizeCard:h,tabTextColorLine:p,tabTextColorActiveLine:n,tabTextColorHoverLine:n,tabTextColorDisabledLine:r,tabTextColorSegment:p,tabTextColorActiveSegment:t,tabTextColorHoverSegment:t,tabTextColorDisabledSegment:r,tabTextColorBar:p,tabTextColorActiveBar:n,tabTextColorHoverBar:n,tabTextColorDisabledBar:r,tabTextColorCard:p,tabTextColorHoverCard:p,tabTextColorActiveCard:n,tabTextColorDisabledCard:r,barColor:n,closeIconColor:i,closeIconColorHover:a,closeIconColorPressed:o,closeColorHover:s,closeColorPressed:c,closeBorderRadius:m,tabColor:l,tabColorSegment:u,tabBorderColor:d,tabFontWeightActive:f,tabFontWeight:f,tabBorderRadius:m,paneTextColor:t,fontWeightStrong:g})}var at={name:`Tabs`,common:U,self:it},ot=A([q(`page-header-header`,`
 margin-bottom: 20px;
 `),q(`page-header`,`
 display: flex;
 align-items: center;
 justify-content: space-between;
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[a(`main`,`
 display: flex;
 flex-wrap: nowrap;
 align-items: center;
 `),a(`back`,`
 display: flex;
 margin-right: 16px;
 font-size: var(--n-back-size);
 cursor: pointer;
 color: var(--n-back-color);
 transition: color .3s var(--n-bezier);
 `,[A(`&:hover`,`color: var(--n-back-color-hover);`),A(`&:active`,`color: var(--n-back-color-pressed);`)]),a(`avatar`,`
 display: flex;
 margin-right: 12px
 `),a(`title`,`
 margin-right: 16px;
 transition: color .3s var(--n-bezier);
 font-size: var(--n-title-font-size);
 font-weight: var(--n-title-font-weight);
 color: var(--n-title-text-color);
 `),a(`subtitle`,`
 font-size: 14px;
 transition: color .3s var(--n-bezier);
 color: var(--n-subtitle-text-color);
 `)]),q(`page-header-content`,`
 font-size: var(--n-font-size);
 `,[A(`&:not(:first-child)`,`margin-top: 20px;`)]),q(`page-header-footer`,`
 font-size: var(--n-font-size);
 `,[A(`&:not(:first-child)`,`margin-top: 20px;`)])]),st=T({name:`PageHeader`,props:Object.assign(Object.assign({},F.props),{title:String,subtitle:String,extra:String,onBack:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n,inlineThemeDisabled:r}=v(e),i=F(`PageHeader`,`-page-header`,ot,nt,e,t),a=J(`PageHeader`,n,t),o=S(()=>{let{self:{titleTextColor:e,subtitleTextColor:t,backColor:n,fontSize:r,titleFontSize:a,backSize:o,titleFontWeight:s,backColorHover:c,backColorPressed:l},common:{cubicBezierEaseInOut:u}}=i.value;return{"--n-title-text-color":e,"--n-title-font-size":a,"--n-title-font-weight":s,"--n-font-size":r,"--n-back-size":o,"--n-subtitle-text-color":t,"--n-back-color":n,"--n-back-color-hover":c,"--n-back-color-pressed":l,"--n-bezier":u}}),s=r?E(`page-header`,void 0,o,e):void 0;return{rtlEnabled:a,mergedClsPrefix:t,cssVars:r?void 0:o,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{onBack:t,title:n,subtitle:r,extra:i,mergedClsPrefix:a,cssVars:o,$slots:s}=this;(e=this.onRender)==null||e.call(this);let{title:c,subtitle:l,extra:u,default:d,header:f,avatar:p,footer:m,back:h}=s,g=t,_=n||c,v=r||l,y=i||u;return D(`div`,{style:o,class:[`${a}-page-header-wrapper`,this.themeClass,this.rtlEnabled&&`${a}-page-header-wrapper--rtl`]},f?D(`div`,{class:`${a}-page-header-header`,key:`breadcrumb`},f()):null,(g||p||_||v||y)&&D(`div`,{class:`${a}-page-header`,key:`header`},D(`div`,{class:`${a}-page-header__main`,key:`back`},g?D(`div`,{class:`${a}-page-header__back`,onClick:t},h?h():D(R,{clsPrefix:a},{default:()=>D($e,null)})):null,p?D(`div`,{class:`${a}-page-header__avatar`},p()):null,_?D(`div`,{class:`${a}-page-header__title`,key:`title`},n||c()):null,v?D(`div`,{class:`${a}-page-header__subtitle`,key:`subtitle`},r||l()):null),y?D(`div`,{class:`${a}-page-header__extra`},i||u()):null),d?D(`div`,{class:`${a}-page-header-content`,key:`content`},d()):null,m?D(`div`,{class:`${a}-page-header-footer`,key:`footer`},m()):null)}}),ct=V(`n-tabs`),lt={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:`if`},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},ut=T({__TAB_PANE__:!0,name:`TabPane`,alias:[`TabPanel`],props:lt,slots:Object,setup(e){let t=n(ct,null);return t||M(`tab-pane`,"`n-tab-pane` must be placed inside `n-tabs`."),{style:t.paneStyleRef,class:t.paneClassRef,mergedClsPrefix:t.mergedClsPrefixRef}},render(){return D(`div`,{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),dt=T({__TAB__:!0,inheritAttrs:!1,name:`Tab`,props:Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},xe(lt,[`displayDirective`])),setup(e){let{mergedClsPrefixRef:t,valueRef:r,typeRef:i,closableRef:a,tabStyleRef:o,addTabStyleRef:s,tabClassRef:c,addTabClassRef:l,tabChangeIdRef:u,onBeforeLeaveRef:d,triggerRef:f,handleAdd:p,activateTab:m,handleClose:h}=n(ct);return{trigger:f,mergedClosable:S(()=>{if(e.internalAddable)return!1;let{closable:t}=e;return t===void 0?a.value:t}),style:o,addStyle:s,tabClass:c,addTabClass:l,clsPrefix:t,value:r,type:i,handleClose(t){t.stopPropagation(),!e.disabled&&h(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){p();return}let{name:t}=e,n=++u.id;if(t!==r.value){let{value:i}=d;i?Promise.resolve(i(e.name,r.value)).then(e=>{e&&u.id===n&&m(t)}):m(t)}}}},render(){let{internalAddable:e,clsPrefix:t,name:n,disabled:r,label:i,tab:a,value:o,mergedClosable:s,trigger:c,$slots:{default:l}}=this,u=i??a;return D(`div`,{class:`${t}-tabs-tab-wrapper`},this.internalLeftPadded?D(`div`,{class:`${t}-tabs-tab-pad`}):null,D(`div`,Object.assign({key:n,"data-name":n,"data-disabled":r?!0:void 0},P({class:[`${t}-tabs-tab`,o===n&&`${t}-tabs-tab--active`,r&&`${t}-tabs-tab--disabled`,s&&`${t}-tabs-tab--closable`,e&&`${t}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:c===`click`?this.activateTab:void 0,onMouseenter:c===`hover`?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),D(`span`,{class:`${t}-tabs-tab__label`},e?D(p,null,D(`div`,{class:`${t}-tabs-tab__height-placeholder`},`\xA0`),D(R,{clsPrefix:t},{default:()=>D(me,null)})):l?l():typeof u==`object`?u:ve(u??n)),s&&this.type===`card`?D(ce,{clsPrefix:t,class:`${t}-tabs-tab__close`,onClick:this.handleClose,disabled:r}):null))}}),ft=q(`tabs`,`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[K(`segment-type`,[q(`tabs-rail`,[A(`&.transition-disabled`,[q(`tabs-capsule`,`
 transition: none;
 `)])])]),K(`top`,[q(`tab-pane`,`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),K(`left`,[q(`tab-pane`,`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),K(`left, right`,`
 flex-direction: row;
 `,[q(`tabs-bar`,`
 width: 2px;
 right: 0;
 transition:
 top .2s var(--n-bezier),
 max-height .2s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),q(`tabs-tab`,`
 padding: var(--n-tab-padding-vertical); 
 `)]),K(`right`,`
 flex-direction: row-reverse;
 `,[q(`tab-pane`,`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),q(`tabs-bar`,`
 left: 0;
 `)]),K(`bottom`,`
 flex-direction: column-reverse;
 justify-content: flex-end;
 `,[q(`tab-pane`,`
 padding: var(--n-pane-padding-bottom) var(--n-pane-padding-right) var(--n-pane-padding-top) var(--n-pane-padding-left);
 `),q(`tabs-bar`,`
 top: 0;
 `)]),q(`tabs-rail`,`
 position: relative;
 padding: 3px;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 background-color: var(--n-color-segment);
 transition: background-color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 `,[q(`tabs-capsule`,`
 border-radius: var(--n-tab-border-radius);
 position: absolute;
 pointer-events: none;
 background-color: var(--n-tab-color-segment);
 box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .08);
 transition: transform 0.3s var(--n-bezier);
 `),q(`tabs-tab-wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[q(`tabs-tab`,`
 overflow: hidden;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[K(`active`,`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),A(`&:hover`,`
 color: var(--n-tab-text-color-hover);
 `)])])]),K(`flex`,[q(`tabs-nav`,`
 width: 100%;
 position: relative;
 `,[q(`tabs-wrapper`,`
 width: 100%;
 `,[q(`tabs-tab`,`
 margin-right: 0;
 `)])])]),q(`tabs-nav`,`
 box-sizing: border-box;
 line-height: 1.5;
 display: flex;
 transition: border-color .3s var(--n-bezier);
 `,[a(`prefix, suffix`,`
 display: flex;
 align-items: center;
 `),a(`prefix`,`padding-right: 16px;`),a(`suffix`,`padding-left: 16px;`)]),K(`top, bottom`,[A(`>`,[q(`tabs-nav`,[q(`tabs-nav-scroll-wrapper`,[A(`&::before`,`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),A(`&::after`,`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),K(`shadow-start`,[A(`&::before`,`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),K(`shadow-end`,[A(`&::after`,`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])])])]),K(`left, right`,[q(`tabs-nav-scroll-content`,`
 flex-direction: column;
 `),A(`>`,[q(`tabs-nav`,[q(`tabs-nav-scroll-wrapper`,[A(`&::before`,`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),A(`&::after`,`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),K(`shadow-start`,[A(`&::before`,`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),K(`shadow-end`,[A(`&::after`,`
 box-shadow: inset 0 -10px 8px -8px rgba(0, 0, 0, .12);
 `)])])])])]),q(`tabs-nav-scroll-wrapper`,`
 flex: 1;
 position: relative;
 overflow: hidden;
 `,[q(`tabs-nav-y-scroll`,`
 height: 100%;
 width: 100%;
 overflow-y: auto; 
 scrollbar-width: none;
 `,[A(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,`
 width: 0;
 height: 0;
 display: none;
 `)]),A(`&::before, &::after`,`
 transition: box-shadow .3s var(--n-bezier);
 pointer-events: none;
 content: "";
 position: absolute;
 z-index: 1;
 `)]),q(`tabs-nav-scroll-content`,`
 display: flex;
 position: relative;
 min-width: 100%;
 min-height: 100%;
 width: fit-content;
 box-sizing: border-box;
 `),q(`tabs-wrapper`,`
 display: inline-flex;
 flex-wrap: nowrap;
 position: relative;
 `),q(`tabs-tab-wrapper`,`
 display: flex;
 flex-wrap: nowrap;
 flex-shrink: 0;
 flex-grow: 0;
 `),q(`tabs-tab`,`
 cursor: pointer;
 white-space: nowrap;
 flex-wrap: nowrap;
 display: inline-flex;
 align-items: center;
 color: var(--n-tab-text-color);
 font-size: var(--n-tab-font-size);
 background-clip: padding-box;
 padding: var(--n-tab-padding);
 transition:
 box-shadow .3s var(--n-bezier),
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[K(`disabled`,{cursor:`not-allowed`}),a(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),a(`label`,`
 display: flex;
 align-items: center;
 z-index: 1;
 `)]),q(`tabs-bar`,`
 position: absolute;
 bottom: 0;
 height: 2px;
 border-radius: 1px;
 background-color: var(--n-bar-color);
 transition:
 left .2s var(--n-bezier),
 max-width .2s var(--n-bezier),
 opacity .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `,[A(`&.transition-disabled`,`
 transition: none;
 `),K(`disabled`,`
 background-color: var(--n-tab-text-color-disabled)
 `)]),q(`tabs-pane-wrapper`,`
 position: relative;
 overflow: hidden;
 transition: max-height .2s var(--n-bezier);
 `),q(`tab-pane`,`
 color: var(--n-pane-text-color);
 width: 100%;
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .2s var(--n-bezier);
 left: 0;
 right: 0;
 top: 0;
 `,[A(`&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),A(`&.next-transition-leave-active, &.prev-transition-leave-active`,`
 position: absolute;
 `),A(`&.next-transition-enter-from, &.prev-transition-leave-to`,`
 transform: translateX(32px);
 opacity: 0;
 `),A(`&.next-transition-leave-to, &.prev-transition-enter-from`,`
 transform: translateX(-32px);
 opacity: 0;
 `),A(`&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to`,`
 transform: translateX(0);
 opacity: 1;
 `)]),q(`tabs-tab-pad`,`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),K(`line-type, bar-type`,[q(`tabs-tab`,`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[A(`&:hover`,{color:`var(--n-tab-text-color-hover)`}),K(`active`,`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),K(`disabled`,{color:`var(--n-tab-text-color-disabled)`})])]),q(`tabs-nav`,[K(`line-type`,[K(`top`,[a(`prefix, suffix`,`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),q(`tabs-nav-scroll-content`,`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),q(`tabs-bar`,`
 bottom: -1px;
 `)]),K(`left`,[a(`prefix, suffix`,`
 border-right: 1px solid var(--n-tab-border-color);
 `),q(`tabs-nav-scroll-content`,`
 border-right: 1px solid var(--n-tab-border-color);
 `),q(`tabs-bar`,`
 right: -1px;
 `)]),K(`right`,[a(`prefix, suffix`,`
 border-left: 1px solid var(--n-tab-border-color);
 `),q(`tabs-nav-scroll-content`,`
 border-left: 1px solid var(--n-tab-border-color);
 `),q(`tabs-bar`,`
 left: -1px;
 `)]),K(`bottom`,[a(`prefix, suffix`,`
 border-top: 1px solid var(--n-tab-border-color);
 `),q(`tabs-nav-scroll-content`,`
 border-top: 1px solid var(--n-tab-border-color);
 `),q(`tabs-bar`,`
 top: -1px;
 `)]),a(`prefix, suffix`,`
 transition: border-color .3s var(--n-bezier);
 `),q(`tabs-nav-scroll-content`,`
 transition: border-color .3s var(--n-bezier);
 `),q(`tabs-bar`,`
 border-radius: 0;
 `)]),K(`card-type`,[a(`prefix, suffix`,`
 transition: border-color .3s var(--n-bezier);
 `),q(`tabs-pad`,`
 flex-grow: 1;
 transition: border-color .3s var(--n-bezier);
 `),q(`tabs-tab-pad`,`
 transition: border-color .3s var(--n-bezier);
 `),q(`tabs-tab`,`
 font-weight: var(--n-tab-font-weight);
 border: 1px solid var(--n-tab-border-color);
 background-color: var(--n-tab-color);
 box-sizing: border-box;
 position: relative;
 vertical-align: bottom;
 display: flex;
 justify-content: space-between;
 font-size: var(--n-tab-font-size);
 color: var(--n-tab-text-color);
 `,[K(`addable`,`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[a(`height-placeholder`,`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),y(`disabled`,[A(`&:hover`,`
 color: var(--n-tab-text-color-hover);
 `)])]),K(`closable`,`padding-right: 8px;`),K(`active`,`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),K(`disabled`,`color: var(--n-tab-text-color-disabled);`)])]),K(`left, right`,`
 flex-direction: column; 
 `,[a(`prefix, suffix`,`
 padding: var(--n-tab-padding-vertical);
 `),q(`tabs-wrapper`,`
 flex-direction: column;
 `),q(`tabs-tab-wrapper`,`
 flex-direction: column;
 `,[q(`tabs-tab-pad`,`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),K(`top`,[K(`card-type`,[q(`tabs-scroll-padding`,`border-bottom: 1px solid var(--n-tab-border-color);`),a(`prefix, suffix`,`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),q(`tabs-tab`,`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[K(`active`,`
 border-bottom: 1px solid #0000;
 `)]),q(`tabs-tab-pad`,`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),q(`tabs-pad`,`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),K(`left`,[K(`card-type`,[q(`tabs-scroll-padding`,`border-right: 1px solid var(--n-tab-border-color);`),a(`prefix, suffix`,`
 border-right: 1px solid var(--n-tab-border-color);
 `),q(`tabs-tab`,`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[K(`active`,`
 border-right: 1px solid #0000;
 `)]),q(`tabs-tab-pad`,`
 border-right: 1px solid var(--n-tab-border-color);
 `),q(`tabs-pad`,`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),K(`right`,[K(`card-type`,[q(`tabs-scroll-padding`,`border-left: 1px solid var(--n-tab-border-color);`),a(`prefix, suffix`,`
 border-left: 1px solid var(--n-tab-border-color);
 `),q(`tabs-tab`,`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[K(`active`,`
 border-left: 1px solid #0000;
 `)]),q(`tabs-tab-pad`,`
 border-left: 1px solid var(--n-tab-border-color);
 `),q(`tabs-pad`,`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),K(`bottom`,[K(`card-type`,[q(`tabs-scroll-padding`,`border-top: 1px solid var(--n-tab-border-color);`),a(`prefix, suffix`,`
 border-top: 1px solid var(--n-tab-border-color);
 `),q(`tabs-tab`,`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[K(`active`,`
 border-top: 1px solid #0000;
 `)]),q(`tabs-tab-pad`,`
 border-top: 1px solid var(--n-tab-border-color);
 `),q(`tabs-pad`,`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),pt=Qe,mt=T({name:`Tabs`,props:Object.assign(Object.assign({},F.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:`click`},type:{type:String,default:`bar`},closable:Boolean,justifyContent:String,size:String,placement:{type:String,default:`top`},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),slots:Object,setup(t,{slots:n}){let{mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:c}=v(t),l=F(`Tabs`,`-tabs`,ft,at,t,i),u=o(null),d=o(null),f=o(null),p=o(null),m=o(null),h=o(null),g=o(!0),_=o(!0),y=ie(t,[`labelSize`,`size`]),b=S(()=>y.value?y.value:c?.value?.Tabs?.size||`medium`),x=ie(t,[`activeName`,`value`]),C=o(x.value??t.defaultValue??(n.default?ye(n.default())[0]?.props?.name:null)),w=e(x,C),T={id:0},D=S(()=>{if(!(!t.justifyContent||t.type===`card`))return{display:`flex`,justifyContent:t.justifyContent}});z(w,()=>{T.id=0,M(),N()});function O(){let{value:e}=w;return e===null?null:u.value?.querySelector(`[data-name="${e}"]`)}function k(e){if(t.type===`card`)return;let{value:n}=d;if(!n)return;let r=n.style.opacity===`0`;if(e){let a=`${i.value}-tabs-bar--disabled`,{barWidth:o,placement:s}=t;if(e.dataset.disabled===`true`?n.classList.add(a):n.classList.remove(a),[`top`,`bottom`].includes(s)){if(j([`top`,`maxHeight`,`height`]),typeof o==`number`&&e.offsetWidth>=o){let t=Math.floor((e.offsetWidth-o)/2)+e.offsetLeft;n.style.left=`${t}px`,n.style.maxWidth=`${o}px`}else n.style.left=`${e.offsetLeft}px`,n.style.maxWidth=`${e.offsetWidth}px`;n.style.width=`8192px`,r&&(n.style.transition=`none`),n.offsetWidth,r&&(n.style.transition=``,n.style.opacity=`1`)}else{if(j([`left`,`maxWidth`,`width`]),typeof o==`number`&&e.offsetHeight>=o){let t=Math.floor((e.offsetHeight-o)/2)+e.offsetTop;n.style.top=`${t}px`,n.style.maxHeight=`${o}px`}else n.style.top=`${e.offsetTop}px`,n.style.maxHeight=`${e.offsetHeight}px`;n.style.height=`8192px`,r&&(n.style.transition=`none`),n.offsetHeight,r&&(n.style.transition=``,n.style.opacity=`1`)}}}function A(){if(t.type===`card`)return;let{value:e}=d;e&&(e.style.opacity=`0`)}function j(e){let{value:t}=d;if(t)for(let n of e)t.style[n]=``}function M(){if(t.type===`card`)return;let e=O();e?k(e):A()}function N(){let e=m.value?.$el;if(!e)return;let t=O();if(!t)return;let{scrollLeft:n,offsetWidth:r}=e,{offsetLeft:i,offsetWidth:a}=t;n>i?e.scrollTo({top:0,left:i,behavior:`smooth`}):i+a>n+r&&e.scrollTo({top:0,left:i+a-r,behavior:`smooth`})}let P=o(null),I=0,L=null;function R(e){let t=P.value;if(t){I=e.getBoundingClientRect().height;let n=`${I}px`,r=()=>{t.style.height=n,t.style.maxHeight=n};L?(r(),L(),L=null):L=r}}function V(e){let t=P.value;if(t){let n=e.getBoundingClientRect().height,r=()=>{document.body.offsetHeight,t.style.maxHeight=`${n}px`,t.style.height=`${Math.max(I,n)}px`};L?(L(),L=null,r()):L=r}}function H(){let e=P.value;if(e){e.style.maxHeight=``,e.style.height=``;let{paneWrapperStyle:n}=t;if(typeof n==`string`)e.style.cssText=n;else if(n){let{maxHeight:t,height:r}=n;t!==void 0&&(e.style.maxHeight=t),r!==void 0&&(e.style.height=r)}}}let U={value:[]},W=o(`next`);function re(e){let t=w.value,n=`next`;for(let r of U.value){if(r===t)break;if(r===e){n=`prev`;break}}W.value=n,K(e)}function K(e){let{onActiveNameChange:n,onUpdateValue:r,"onUpdate:value":i}=t;n&&B(n,e),r&&B(r,e),i&&B(i,e),C.value=e}function q(e){let{onClose:n}=t;n&&B(n,e)}let J=!0;function ae(){let{value:e}=d;if(!e)return;J||=!1;let t=`transition-disabled`;e.classList.add(t),M(),e.classList.remove(t)}let Y=o(null);function X({transitionDisabled:e}){let t=u.value;if(!t)return;e&&t.classList.add(`transition-disabled`);let n=O();n&&Y.value&&(Y.value.style.width=`${n.offsetWidth}px`,Y.value.style.height=`${n.offsetHeight}px`,Y.value.style.transform=`translateX(${n.offsetLeft-le(getComputedStyle(t).paddingLeft)}px)`,e&&Y.value.offsetWidth),e&&t.classList.remove(`transition-disabled`)}z([w],()=>{t.type===`segment`&&G(()=>{X({transitionDisabled:!1})})}),ne(()=>{t.type===`segment`&&X({transitionDisabled:!0})});let ce=0;function ue(e){if(e.contentRect.width===0&&e.contentRect.height===0||ce===e.contentRect.width)return;ce=e.contentRect.width;let{type:n}=t;if((n===`line`||n===`bar`)&&(J||t.justifyContent?.startsWith(`space`))&&ae(),n!==`segment`){let{placement:e}=t;me((e===`top`||e===`bottom`?m.value?.$el:h.value)||null)}}let de=pt(ue,64);z([()=>t.justifyContent,()=>t.size],()=>{G(()=>{let{type:e}=t;(e===`line`||e===`bar`)&&ae()})});let Z=o(!1);function fe(e){let{target:n,contentRect:{width:r,height:i}}=e,a=n.parentElement.parentElement.offsetWidth,o=n.parentElement.parentElement.offsetHeight,{placement:s}=t;if(!Z.value)s===`top`||s===`bottom`?a<r&&(Z.value=!0):o<i&&(Z.value=!0);else{let{value:e}=p;if(!e)return;s===`top`||s===`bottom`?a-r>e.$el.offsetWidth&&(Z.value=!1):o-i>e.$el.offsetHeight&&(Z.value=!1)}me(m.value?.$el||null)}let pe=pt(fe,64);function Q(){let{onAdd:e}=t;e&&e(),G(()=>{let e=O(),{value:t}=m;!e||!t||t.scrollTo({left:e.offsetLeft,top:0,behavior:`smooth`})})}function me(e){if(!e)return;let{placement:n}=t;if(n===`top`||n===`bottom`){let{scrollLeft:t,scrollWidth:n,offsetWidth:r}=e;g.value=t<=0,_.value=t+r>=n}else{let{scrollTop:t,scrollHeight:n,offsetHeight:r}=e;g.value=t<=0,_.value=t+r>=n}}let he=pt(e=>{me(e.target)},64);te(ct,{triggerRef:r(t,`trigger`),tabStyleRef:r(t,`tabStyle`),tabClassRef:r(t,`tabClass`),addTabStyleRef:r(t,`addTabStyle`),addTabClassRef:r(t,`addTabClass`),paneClassRef:r(t,`paneClass`),paneStyleRef:r(t,`paneStyle`),mergedClsPrefixRef:i,typeRef:r(t,`type`),closableRef:r(t,`closable`),valueRef:w,tabChangeIdRef:T,onBeforeLeaveRef:r(t,`onBeforeLeave`),activateTab:re,handleClose:q,handleAdd:Q}),oe(()=>{M(),N()}),ee(()=>{let{value:e}=f;if(!e)return;let{value:t}=i,n=`${t}-tabs-nav-scroll-wrapper--shadow-start`,r=`${t}-tabs-nav-scroll-wrapper--shadow-end`;g.value?e.classList.remove(n):e.classList.add(n),_.value?e.classList.remove(r):e.classList.add(r)});let ge={syncBarPosition:()=>{M()}},_e=()=>{X({transitionDisabled:!0})},$=S(()=>{let{value:e}=b,{type:n}=t,r=`${e}${{card:`Card`,bar:`Bar`,line:`Line`,segment:`Segment`}[n]}`,{self:{barColor:i,closeIconColor:a,closeIconColorHover:o,closeIconColorPressed:c,tabColor:u,tabBorderColor:d,paneTextColor:f,tabFontWeight:p,tabBorderRadius:m,tabFontWeightActive:h,colorSegment:g,fontWeightStrong:_,tabColorSegment:v,closeSize:y,closeIconSize:x,closeColorHover:S,closeColorPressed:C,closeBorderRadius:w,[s(`panePadding`,e)]:T,[s(`tabPadding`,r)]:E,[s(`tabPaddingVertical`,r)]:D,[s(`tabGap`,r)]:O,[s(`tabGap`,`${r}Vertical`)]:k,[s(`tabTextColor`,n)]:A,[s(`tabTextColorActive`,n)]:j,[s(`tabTextColorHover`,n)]:M,[s(`tabTextColorDisabled`,n)]:N,[s(`tabFontSize`,e)]:P},common:{cubicBezierEaseInOut:F}}=l.value;return{"--n-bezier":F,"--n-color-segment":g,"--n-bar-color":i,"--n-tab-font-size":P,"--n-tab-text-color":A,"--n-tab-text-color-active":j,"--n-tab-text-color-disabled":N,"--n-tab-text-color-hover":M,"--n-pane-text-color":f,"--n-tab-border-color":d,"--n-tab-border-radius":m,"--n-close-size":y,"--n-close-icon-size":x,"--n-close-color-hover":S,"--n-close-color-pressed":C,"--n-close-border-radius":w,"--n-close-icon-color":a,"--n-close-icon-color-hover":o,"--n-close-icon-color-pressed":c,"--n-tab-color":u,"--n-tab-font-weight":p,"--n-tab-font-weight-active":h,"--n-tab-padding":E,"--n-tab-padding-vertical":D,"--n-tab-gap":O,"--n-tab-gap-vertical":k,"--n-pane-padding-left":se(T,`left`),"--n-pane-padding-right":se(T,`right`),"--n-pane-padding-top":se(T,`top`),"--n-pane-padding-bottom":se(T,`bottom`),"--n-font-weight-strong":_,"--n-tab-color-segment":v}}),ve=a?E(`tabs`,S(()=>`${b.value[0]}${t.type[0]}`),$,t):void 0;return Object.assign({mergedClsPrefix:i,mergedValue:w,renderedNames:new Set,segmentCapsuleElRef:Y,tabsPaneWrapperRef:P,tabsElRef:u,barElRef:d,addTabInstRef:p,xScrollInstRef:m,scrollWrapperElRef:f,addTabFixed:Z,tabWrapperStyle:D,handleNavResize:de,mergedSize:b,handleScroll:he,handleTabsResize:pe,cssVars:a?void 0:$,themeClass:ve?.themeClass,animationDirection:W,renderNameListRef:U,yScrollElRef:h,handleSegmentResize:_e,onAnimationBeforeLeave:R,onAnimationEnter:V,onAnimationAfterEnter:H,onRender:ve?.onRender},ge)},render(){let{mergedClsPrefix:e,type:t,placement:n,addTabFixed:r,addable:i,mergedSize:a,renderNameListRef:o,onRender:s,paneWrapperClass:c,paneWrapperStyle:l,$slots:{default:u,prefix:d,suffix:f}}=this;s?.();let p=u?ye(u()).filter(e=>e.type.__TAB_PANE__===!0):[],m=u?ye(u()).filter(e=>e.type.__TAB__===!0):[],h=!m.length,g=t===`card`,_=t===`segment`,v=!g&&!_&&this.justifyContent;o.value=[];let y=()=>{let t=D(`div`,{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},v?null:D(`div`,{class:`${e}-tabs-scroll-padding`,style:n===`top`||n===`bottom`?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),h?p.map((e,t)=>(o.value.push(e.props.name),vt(D(dt,Object.assign({},e.props,{internalCreatedByPane:!0,internalLeftPadded:t!==0&&(!v||v===`center`||v===`start`||v===`end`)}),e.children?{default:e.children.tab}:void 0)))):m.map((e,t)=>(o.value.push(e.props.name),vt(t!==0&&!v?_t(e):e))),!r&&i&&g?gt(i,(h?p.length:m.length)!==0):null,v?null:D(`div`,{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return D(`div`,{ref:`tabsElRef`,class:`${e}-tabs-nav-scroll-content`},g&&i?D(X,{onResize:this.handleTabsResize},{default:()=>t}):t,g?D(`div`,{class:`${e}-tabs-pad`}):null,g?null:D(`div`,{ref:`barElRef`,class:`${e}-tabs-bar`}))},b=_?`top`:n;return D(`div`,{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${t}-type`,`${e}-tabs--${a}-size`,v&&`${e}-tabs--flex`,`${e}-tabs--${b}`],style:this.cssVars},D(`div`,{class:[`${e}-tabs-nav--${t}-type`,`${e}-tabs-nav--${b}`,`${e}-tabs-nav`]},k(d,t=>t&&D(`div`,{class:`${e}-tabs-nav__prefix`},t)),_?D(X,{onResize:this.handleSegmentResize},{default:()=>D(`div`,{class:`${e}-tabs-rail`,ref:`tabsElRef`},D(`div`,{class:`${e}-tabs-capsule`,ref:`segmentCapsuleElRef`},D(`div`,{class:`${e}-tabs-wrapper`},D(`div`,{class:`${e}-tabs-tab`}))),h?p.map((e,t)=>(o.value.push(e.props.name),D(dt,Object.assign({},e.props,{internalCreatedByPane:!0,internalLeftPadded:t!==0}),e.children?{default:e.children.tab}:void 0))):m.map((e,t)=>(o.value.push(e.props.name),t===0?e:_t(e))))}):D(X,{onResize:this.handleNavResize},{default:()=>D(`div`,{class:`${e}-tabs-nav-scroll-wrapper`,ref:`scrollWrapperElRef`},[`top`,`bottom`].includes(b)?D(Fe,{ref:`xScrollInstRef`,onScroll:this.handleScroll},{default:y}):D(`div`,{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:`yScrollElRef`},y()))}),r&&i&&g?gt(i,!0):null,k(f,t=>t&&D(`div`,{class:`${e}-tabs-nav__suffix`},t))),h&&(this.animated&&(b===`top`||b===`bottom`)?D(`div`,{ref:`tabsPaneWrapperRef`,style:l,class:[`${e}-tabs-pane-wrapper`,c]},ht(p,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):ht(p,this.mergedValue,this.renderedNames)))}});function ht(e,t,n,r,i,a,o){let s=[];return e.forEach(e=>{let{name:r,displayDirective:i,"display-directive":a}=e.props,o=e=>i===e||a===e,c=t===r;if(e.key!==void 0&&(e.key=r),c||o(`show`)||o(`show:lazy`)&&n.has(r)){n.has(r)||n.add(r);let t=!o(`if`);s.push(t?O(e,[[_,c]]):e)}}),o?D(h,{name:`${o}-transition`,onBeforeLeave:r,onEnter:i,onAfterEnter:a},{default:()=>s}):s}function gt(e,t){return D(dt,{ref:`addTabInstRef`,key:`__addable`,name:`__addable`,internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:t,disabled:typeof e==`object`&&e.disabled})}function _t(e){let t=i(e);return t.props?t.props.internalLeftPadded=!0:t.props={internalLeftPadded:!0},t}function vt(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes(`internalLeftPadded`)||e.dynamicProps.push(`internalLeftPadded`):e.dynamicProps=[`internalLeftPadded`],e}var yt={class:`device-detail`},bt={class:`config-content`},xt={key:0},St={class:`diff-viewer`},Ct={key:1},wt={class:`diff-content`},Tt=we(T({__name:`DeviceDetail`,setup(e){let n=be(),r=Ce(),i=t(),a=o(!1),s=o(!1),d=o(!1),h=o(!1),_=o(!1),v=o(!1),y=o(null),x=o(``),C=o([]),T=o([]),E=o(`latest`),D=Number(r.params.id),O=Number.isNaN(D)||D<=0,k=o(null),A=o(null),M=o(null),N=o(!1),P={cisco:`思科`,h3c:`华三`,huawei:`华为`,ruijie:`锐捷`},F={ios:`IOS`,vrp:`VRP`,comware:`Comware`,"rg-os":`RG-OS`},L=S(()=>y.value?P[y.value.vendor]||y.value.vendor:``),ee=S(()=>y.value?F[y.value.model]||y.value.model:``),te=S(()=>C.value.map(e=>({label:`${e.hash.substring(0,8)} - ${e.message} (${e.date})`,value:e.hash}))),R=S(()=>M.value?.diff_content?M.value.diff_content.split(`
`):[]),B=[{title:`Hash`,key:`hash`,ellipsis:{tooltip:!0}},{title:`提交信息`,key:`message`},{title:`作者`,key:`author`,width:120},{title:`日期`,key:`date`,width:180}];async function V(){a.value=!0;try{y.value=(await Ee(D)).data}catch{i.error(`获取设备信息失败`)}finally{a.value=!1}}async function U(){s.value=!0;try{x.value=(await Me(D)).data.content}catch{i.error(`获取设备配置失败`)}finally{s.value=!1}}async function G(){d.value=!0;try{C.value=(await Ae(D)).data}catch{i.error(`获取历史版本失败`)}finally{d.value=!1}}async function K(){h.value=!0;try{T.value=(await je(D,{page:1,page_size:20})).data.items}catch{i.error(`获取基线偏差失败`)}finally{h.value=!1}}async function q(){if(!(!k.value||!A.value)){v.value=!0;try{M.value=(await Ne(D,k.value,A.value)).data,N.value=!0}catch{i.error(`配置比对失败`)}finally{v.value=!1}}}async function J(){_.value=!0;try{await ke(D),i.success(`备份任务已触发`)}catch{i.error(`备份任务触发失败`)}finally{_.value=!1}}async function ie(){if(y.value)try{y.value.enabled?(await De(y.value.id),i.success(`设备已禁用`)):(await Te(y.value.id),i.success(`设备已启用`)),V()}catch{i.error(y.value.enabled?`禁用设备失败`:`启用设备失败`)}}return ne(()=>{if(O){i.error(`无效的设备 ID`),n.push(`/devices`);return}V(),U()}),z(E,e=>{e===`latest`?U():e===`history`?G():e===`deviations`&&K()}),(e,t)=>(j(),l(`div`,yt,[f(m(st),{onBack:t[1]||=e=>m(n).push(`/devices`),title:`设备详情`,subtitle:``},{extra:I(()=>[f(m(he),null,{default:I(()=>[f(m(W),{onClick:t[0]||=e=>m(n).push(`/devices`)},{default:I(()=>[...t[7]||=[H(`返回列表`,-1)]]),_:1}),f(m(W),{type:y.value?.enabled?`warning`:`success`,onClick:ie,disabled:!y.value},{default:I(()=>[H(c(y.value?.enabled?`禁用设备`:`启用设备`),1)]),_:1},8,[`type`,`disabled`]),f(m(W),{type:`primary`,loading:_.value,onClick:J},{default:I(()=>[...t[8]||=[H(`触发备份`,-1)]]),_:1},8,[`loading`])]),_:1})]),_:1}),f(m(ge),{show:a.value},{default:I(()=>[f(m($),{title:`设备信息`,style:{"margin-top":`16px`}},{default:I(()=>[f(m(pe),{bordered:``,column:3},{default:I(()=>[f(m(Q),{label:`设备名称`},{default:I(()=>[H(c(y.value?.name),1)]),_:1}),f(m(Q),{label:`IP 地址`},{default:I(()=>[H(c(y.value?.ip),1)]),_:1}),f(m(Q),{label:`厂商`},{default:I(()=>[H(c(L.value),1)]),_:1}),f(m(Q),{label:`型号`},{default:I(()=>[H(c(ee.value),1)]),_:1}),f(m(Q),{label:`连接协议`},{default:I(()=>[H(c(y.value?.protocol?.toUpperCase()),1)]),_:1}),f(m(Q),{label:`端口`},{default:I(()=>[H(c(y.value?.port),1)]),_:1}),f(m(Q),{label:`分组`},{default:I(()=>[H(c(y.value?.group_name||`-`),1)]),_:1}),f(m(Q),{label:`启用状态`},{default:I(()=>[f(m(ue),{type:y.value?.enabled?`success`:`default`,size:`small`},{default:I(()=>[H(c(y.value?.enabled?`已启用`:`已禁用`),1)]),_:1},8,[`type`])]),_:1}),f(m(Q),{label:`备份状态`},{default:I(()=>[y.value?(j(),b(Oe,{key:0,status:y.value.last_backup_status},null,8,[`status`])):w(``,!0)]),_:1}),f(m(Q),{label:`最近备份`},{default:I(()=>[H(c(y.value?.last_backup_at||`-`),1)]),_:1}),f(m(Q),{label:`失败原因`,span:2},{default:I(()=>[H(c(y.value?.last_failure_reason||`-`),1)]),_:1})]),_:1})]),_:1}),f(m($),{style:{"margin-top":`16px`}},{default:I(()=>[f(m(mt),{value:E.value,"onUpdate:value":t[6]||=e=>E.value=e,type:`line`},{default:I(()=>[f(m(ut),{name:`latest`,tab:`最新配置`},{default:I(()=>[f(m(ge),{show:s.value},{default:I(()=>[u(`div`,bt,[u(`pre`,null,[u(`code`,null,c(x.value||`暂无配置`),1)])])]),_:1},8,[`show`])]),_:1}),f(m(ut),{name:`history`,tab:`历史版本`},{default:I(()=>[C.value.length>=2?(j(),b(m(he),{key:0,style:{"margin-bottom":`12px`}},{default:I(()=>[f(m(Z),{value:k.value,"onUpdate:value":t[2]||=e=>k.value=e,options:te.value,placeholder:`选择源版本`,style:{width:`300px`}},null,8,[`value`,`options`]),f(m(Z),{value:A.value,"onUpdate:value":t[3]||=e=>A.value=e,options:te.value,placeholder:`选择目标版本`,style:{width:`300px`}},null,8,[`value`,`options`]),f(m(W),{type:`primary`,disabled:!k.value||!A.value||k.value===A.value,loading:v.value,onClick:q},{default:I(()=>[...t[9]||=[H(` 比对 `,-1)]]),_:1},8,[`disabled`,`loading`])]),_:1})):w(``,!0),f(m(fe),{columns:B,data:C.value,loading:d.value,size:`small`},null,8,[`data`,`loading`]),f(m(Se),{show:N.value,"onUpdate:show":t[5]||=e=>N.value=e,preset:`dialog`,title:`配置比对结果`,style:{width:`800px`}},{action:I(()=>[f(m(W),{onClick:t[4]||=e=>N.value=!1},{default:I(()=>[...t[10]||=[H(`关闭`,-1)]]),_:1})]),default:I(()=>[M.value?(j(),l(`div`,xt,[f(m(he),{style:{"margin-bottom":`8px`}},{default:I(()=>[f(m(ue),{type:`success`},{default:I(()=>[H(`新增 `+c(M.value.added_lines)+` 行`,1)]),_:1}),f(m(ue),{type:`error`},{default:I(()=>[H(`删除 `+c(M.value.removed_lines)+` 行`,1)]),_:1})]),_:1}),u(`div`,St,[(j(!0),l(p,null,re(R.value,(e,t)=>(j(),l(`pre`,{key:t,class:g({"diff-line-added":e.startsWith(`+`),"diff-line-removed":e.startsWith(`-`),"diff-line-context":e.startsWith(` `)})},c(e),3))),128))])])):w(``,!0)]),_:1},8,[`show`])]),_:1}),f(m(ut),{name:`deviations`,tab:`基线偏差`},{default:I(()=>[f(m(ge),{show:h.value},{default:I(()=>[!h.value&&T.value.length===0?(j(),b(m(de),{key:0,description:`暂无偏差`})):(j(),l(`div`,Ct,[(j(!0),l(p,null,re(T.value,e=>(j(),b(m($),{key:e.id,size:`small`,style:{"margin-bottom":`8px`}},{footer:I(()=>[f(m(_e),{depth:`3`},{default:I(()=>[H(`检测时间: `+c(e.detected_at),1)]),_:2},1024)]),default:I(()=>[u(`pre`,wt,c(e.diff_content),1)]),_:2},1024))),128))]))]),_:1},8,[`show`])]),_:1})]),_:1},8,[`value`])]),_:1})]),_:1},8,[`show`])]))}}),[[`__scopeId`,`data-v-1d7bf808`]]);export{Tt as default};