import{m as e,s as t}from"./api-CZuskCjY.js";import{$t as n,An as r,Ct as i,En as a,Et as o,Lt as s,Q as c,Tn as l,Tt as u,Vt as d,Yt as f,Z as p,Zt as m,bt as h,en as g,g as _,gt as v,h as y,ht as b,in as x,jt as S,ln as C,lt as w,mn as T,ot as E,pt as D,rn as O,s as k,wn as A,wt as j,xt as M}from"./Button-v7qKim5p.js";import{S as N,a as P,f as F,g as I,h as L,i as R,o as z,r as B,s as ee}from"./Tooltip-qMw84P-o.js";import{b as V,d as te,x as H}from"./token-dQuugHoF.js";import{D as ne,O as re,f as U,j as ie,k as ae,m as oe,s as se}from"./index-CESIsNoN.js";function W(e={},t){let n=A({ctrl:!1,command:!1,win:!1,shift:!1,tab:!1}),{keydown:r,keyup:i}=e,a=e=>{switch(e.key){case`Control`:n.ctrl=!0;break;case`Meta`:n.command=!0,n.win=!0;break;case`Shift`:n.shift=!0;break;case`Tab`:n.tab=!0;break}r!==void 0&&Object.keys(r).forEach(t=>{if(t!==e.key)return;let n=r[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},o=e=>{switch(e.key){case`Control`:n.ctrl=!1;break;case`Meta`:n.command=!1,n.win=!1;break;case`Shift`:n.shift=!1;break;case`Tab`:n.tab=!1;break}i!==void 0&&Object.keys(i).forEach(t=>{if(t!==e.key)return;let n=i[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},s=()=>{(t===void 0||t.value)&&(H(`keydown`,document,a),H(`keyup`,document,o)),t!==void 0&&T(t,e=>{e?(H(`keydown`,document,a),H(`keyup`,document,o)):(V(`keydown`,document,a),V(`keyup`,document,o))})};return ie()?(O(s),x(()=>{(t===void 0||t.value)&&(V(`keydown`,document,a),V(`keyup`,document,o))})):s(),l(n)}function ce(e,t,n){if(!t)return e;let r=a(e.value),i=null;return T(e,e=>{i!==null&&window.clearTimeout(i),e===!0?n&&!n.value?r.value=!0:i=window.setTimeout(()=>{r.value=!0},t):r.value=!1}),r}function G(e){return t=>{t?e.value=t.$el:e.value=null}}var K=f({name:`ChevronRight`,render(){return m(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},m(`path`,{d:`M5.64645 3.14645C5.45118 3.34171 5.45118 3.65829 5.64645 3.85355L9.79289 8L5.64645 12.1464C5.45118 12.3417 5.45118 12.6583 5.64645 12.8536C5.84171 13.0488 6.15829 13.0488 6.35355 12.8536L10.8536 8.35355C11.0488 8.15829 11.0488 7.84171 10.8536 7.64645L6.35355 3.14645C6.15829 2.95118 5.84171 2.95118 5.64645 3.14645Z`,fill:`currentColor`}))}}),le={padding:`4px 0`,optionIconSizeSmall:`14px`,optionIconSizeMedium:`16px`,optionIconSizeLarge:`16px`,optionIconSizeHuge:`18px`,optionSuffixWidthSmall:`14px`,optionSuffixWidthMedium:`14px`,optionSuffixWidthLarge:`16px`,optionSuffixWidthHuge:`16px`,optionIconSuffixWidthSmall:`32px`,optionIconSuffixWidthMedium:`32px`,optionIconSuffixWidthLarge:`36px`,optionIconSuffixWidthHuge:`36px`,optionPrefixWidthSmall:`14px`,optionPrefixWidthMedium:`14px`,optionPrefixWidthLarge:`16px`,optionPrefixWidthHuge:`16px`,optionIconPrefixWidthSmall:`36px`,optionIconPrefixWidthMedium:`36px`,optionIconPrefixWidthLarge:`40px`,optionIconPrefixWidthHuge:`40px`};function ue(e){let{primaryColor:t,textColor2:n,dividerColor:r,hoverColor:i,popoverColor:a,invertedColor:o,borderRadius:s,fontSizeSmall:c,fontSizeMedium:l,fontSizeLarge:u,fontSizeHuge:d,heightSmall:f,heightMedium:p,heightLarge:m,heightHuge:h,textColor3:g,opacityDisabled:_}=e;return Object.assign(Object.assign({},le),{optionHeightSmall:f,optionHeightMedium:p,optionHeightLarge:m,optionHeightHuge:h,borderRadius:s,fontSizeSmall:c,fontSizeMedium:l,fontSizeLarge:u,fontSizeHuge:d,optionTextColor:n,optionTextColorHover:n,optionTextColorActive:t,optionTextColorChildActive:t,color:a,dividerColor:r,suffixColor:n,prefixColor:n,optionColorHover:i,optionColorActive:v(t,{alpha:.1}),groupHeaderTextColor:g,optionTextColorInverted:`#BBB`,optionTextColorHoverInverted:`#FFF`,optionTextColorActiveInverted:`#FFF`,optionTextColorChildActiveInverted:`#FFF`,colorInverted:o,dividerColorInverted:`#BBB`,suffixColorInverted:`#BBB`,prefixColorInverted:`#BBB`,optionColorHoverInverted:t,optionColorActiveInverted:t,groupHeaderTextColorInverted:`#AAA`,optionOpacityDisabled:_})}var q=y({name:`Dropdown`,common:k,peers:{Popover:z},self:ue}),J=D(`n-dropdown-menu`),Y=D(`n-dropdown`),X=D(`n-dropdown-option`),Z=f({name:`DropdownDivider`,props:{clsPrefix:{type:String,required:!0}},render(){return m(`div`,{class:`${this.clsPrefix}-dropdown-divider`})}}),de=f({name:`DropdownGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{showIconRef:e,hasSubmenuRef:t}=n(J),{renderLabelRef:r,labelFieldRef:i,nodePropsRef:a,renderOptionRef:o}=n(Y);return{labelField:i,showIcon:e,hasSubmenu:t,renderLabel:r,nodeProps:a,renderOption:o}},render(){let{clsPrefix:e,hasSubmenu:t,showIcon:n,nodeProps:r,renderLabel:i,renderOption:a}=this,{rawNode:o}=this.tmNode,s=m(`div`,Object.assign({class:`${e}-dropdown-option`},r?.(o)),m(`div`,{class:`${e}-dropdown-option-body ${e}-dropdown-option-body--group`},m(`div`,{"data-dropdown-option":!0,class:[`${e}-dropdown-option-body__prefix`,n&&`${e}-dropdown-option-body__prefix--show-icon`]},U(o.icon)),m(`div`,{class:`${e}-dropdown-option-body__label`,"data-dropdown-option":!0},i?i(o):U(o.title??o[this.labelField])),m(`div`,{class:[`${e}-dropdown-option-body__suffix`,t&&`${e}-dropdown-option-body__suffix--has-submenu`],"data-dropdown-option":!0})));return a?a({node:s,option:o}):s}});function Q(e,t){return e.type===`submenu`||e.type===void 0&&e[t]!==void 0}function fe(e){return e.type===`group`}function $(e){return e.type===`divider`}function pe(e){return e.type===`render`}var me=f({name:`DropdownOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null},placement:{type:String,default:`right-start`},props:Object,scrollable:Boolean},setup(e){let t=n(Y),{hoverKeyRef:r,keyboardKeyRef:i,lastToggledSubmenuKeyRef:o,pendingKeyPathRef:s,activeKeyPathRef:c,animatedRef:l,mergedShowRef:u,renderLabelRef:f,renderIconRef:p,labelFieldRef:m,childrenFieldRef:h,renderOptionRef:g,nodePropsRef:_,menuPropsRef:v}=t,y=n(X,null),x=n(J),S=n(ne),w=d(()=>e.tmNode.rawNode),T=d(()=>{let{value:t}=h;return Q(e.tmNode.rawNode,t)}),E=d(()=>{let{disabled:t}=e.tmNode;return t}),D=ce(d(()=>{if(!T.value)return!1;let{key:t,disabled:n}=e.tmNode;if(n)return!1;let{value:a}=r,{value:c}=i,{value:l}=o,{value:u}=s;return a===null?c===null?l===null?!1:u.includes(t):u.includes(t)&&u[u.length-1]!==t:u.includes(t)}),300,d(()=>i.value===null&&!l.value)),O=d(()=>!!y?.enteringSubmenuRef.value),k=a(!1);C(X,{enteringSubmenuRef:k});function A(){k.value=!0}function j(){k.value=!1}function M(){let{parentKey:t,tmNode:n}=e;n.disabled||u.value&&(o.value=t,i.value=null,r.value=n.key)}function P(){let{tmNode:t}=e;t.disabled||u.value&&r.value!==t.key&&M()}function F(t){if(e.tmNode.disabled||!u.value)return;let{relatedTarget:n}=t;n&&!N({target:n},`dropdownOption`)&&!N({target:n},`scrollbarRail`)&&(r.value=null)}function I(){let{value:n}=T,{tmNode:r}=e;u.value&&!n&&!r.disabled&&(t.doSelect(r.key,r.rawNode),t.doUpdateShow(!1))}return{labelField:m,renderLabel:f,renderIcon:p,siblingHasIcon:x.showIconRef,siblingHasSubmenu:x.hasSubmenuRef,menuProps:v,popoverBody:S,animated:l,mergedShowSubmenu:d(()=>D.value&&!O.value),rawNode:w,hasSubmenu:T,pending:b(()=>{let{value:t}=s,{key:n}=e.tmNode;return t.includes(n)}),childActive:b(()=>{let{value:t}=c,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r===-1?!1:r<t.length-1}),active:b(()=>{let{value:t}=c,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r===-1?!1:r===t.length-1}),mergedDisabled:E,renderOption:g,nodeProps:_,handleClick:I,handleMouseMove:P,handleMouseEnter:M,handleMouseLeave:F,handleSubmenuBeforeEnter:A,handleSubmenuAfterEnter:j}},render(){let{animated:e,rawNode:n,mergedShowSubmenu:r,clsPrefix:i,siblingHasIcon:a,siblingHasSubmenu:o,renderLabel:s,renderIcon:c,renderOption:l,nodeProps:u,props:d,scrollable:f}=this,p=null;if(r){let e=this.menuProps?.call(this,n,n.children);p=m(_e,Object.assign({},e,{clsPrefix:i,scrollable:this.scrollable,tmNodes:this.tmNode.children,parentKey:this.tmNode.key}))}let h={class:[`${i}-dropdown-option-body`,this.pending&&`${i}-dropdown-option-body--pending`,this.active&&`${i}-dropdown-option-body--active`,this.childActive&&`${i}-dropdown-option-body--child-active`,this.mergedDisabled&&`${i}-dropdown-option-body--disabled`],onMousemove:this.handleMouseMove,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onClick:this.handleClick},_=u?.(n),v=m(`div`,Object.assign({class:[`${i}-dropdown-option`,_?.class],"data-dropdown-option":!0},_),m(`div`,g(h,d),[m(`div`,{class:[`${i}-dropdown-option-body__prefix`,a&&`${i}-dropdown-option-body__prefix--show-icon`]},[c?c(n):U(n.icon)]),m(`div`,{"data-dropdown-option":!0,class:`${i}-dropdown-option-body__label`},s?s(n):U(n[this.labelField]??n.title)),m(`div`,{"data-dropdown-option":!0,class:[`${i}-dropdown-option-body__suffix`,o&&`${i}-dropdown-option-body__suffix--has-submenu`]},this.hasSubmenu?m(t,null,{default:()=>m(K,null)}):null)]),this.hasSubmenu?m(I,null,{default:()=>[m(L,null,{default:()=>m(`div`,{class:`${i}-dropdown-offset-container`},m(F,{show:this.mergedShowSubmenu,placement:this.placement,to:f&&this.popoverBody||void 0,teleportDisabled:!f},{default:()=>m(`div`,{class:`${i}-dropdown-menu-wrapper`},e?m(S,{onBeforeEnter:this.handleSubmenuBeforeEnter,onAfterEnter:this.handleSubmenuAfterEnter,name:`fade-in-scale-up-transition`,appear:!0},{default:()=>p}):p)}))})]}):null);return l?l({node:v,option:n}):v}}),he=f({name:`NDropdownGroup`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null}},render(){let{tmNode:e,parentKey:t,clsPrefix:n}=this,{children:r}=e;return m(s,null,m(de,{clsPrefix:n,tmNode:e,key:e.key}),r?.map(e=>{let{rawNode:r}=e;return r.show===!1?null:$(r)?m(Z,{clsPrefix:n,key:e.key}):e.isGroup?(w(`dropdown`,"`group` node is not allowed to be put in `group` node."),null):m(me,{clsPrefix:n,tmNode:e,parentKey:t,key:e.key})}))}}),ge=f({name:`DropdownRenderOption`,props:{tmNode:{type:Object,required:!0}},render(){let{rawNode:{render:e,props:t}}=this.tmNode;return m(`div`,t,[e?.()])}}),_e=f({name:`DropdownMenu`,props:{scrollable:Boolean,showArrow:Boolean,arrowStyle:[String,Object],clsPrefix:{type:String,required:!0},tmNodes:{type:Array,default:()=>[]},parentKey:{type:[String,Number],default:null}},setup(e){let{renderIconRef:t,childrenFieldRef:r}=n(Y);C(J,{showIconRef:d(()=>{let n=t.value;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>n?n(e):e.icon);let{rawNode:t}=e;return n?n(t):t.icon})}),hasSubmenuRef:d(()=>{let{value:t}=r;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>Q(e,t));let{rawNode:n}=e;return Q(n,t)})})});let i=a(null);return C(re,null),C(ae,null),C(ne,i),{bodyRef:i}},render(){let{parentKey:e,clsPrefix:t,scrollable:n}=this,r=this.tmNodes.map(r=>{let{rawNode:i}=r;return i.show===!1?null:pe(i)?m(ge,{tmNode:r,key:r.key}):$(i)?m(Z,{clsPrefix:t,key:r.key}):fe(i)?m(he,{clsPrefix:t,tmNode:r,parentKey:e,key:r.key}):m(me,{clsPrefix:t,tmNode:r,parentKey:e,key:r.key,props:i.props,scrollable:n})});return m(`div`,{class:[`${t}-dropdown-menu`,n&&`${t}-dropdown-menu--scrollable`],ref:`bodyRef`},n?m(te,{contentClass:`${t}-dropdown-menu__content`},{default:()=>r}):r,this.showArrow?P({clsPrefix:t,arrowStyle:this.arrowStyle,arrowClass:void 0,arrowWrapperClass:void 0,arrowWrapperStyle:void 0}):null)}}),ve=M(`dropdown-menu`,`
 transform-origin: var(--v-transform-origin);
 background-color: var(--n-color);
 border-radius: var(--n-border-radius);
 box-shadow: var(--n-box-shadow);
 position: relative;
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
`,[se(),M(`dropdown-option`,`
 position: relative;
 `,[h(`a`,`
 text-decoration: none;
 color: inherit;
 outline: none;
 `,[h(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),M(`dropdown-option-body`,`
 display: flex;
 cursor: pointer;
 position: relative;
 height: var(--n-option-height);
 line-height: var(--n-option-height);
 font-size: var(--n-font-size);
 color: var(--n-option-text-color);
 transition: color .3s var(--n-bezier);
 `,[h(`&::before`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 left: 4px;
 right: 4px;
 transition: background-color .3s var(--n-bezier);
 border-radius: var(--n-border-radius);
 `),u(`disabled`,[j(`pending`,`
 color: var(--n-option-text-color-hover);
 `,[i(`prefix, suffix`,`
 color: var(--n-option-text-color-hover);
 `),h(`&::before`,`background-color: var(--n-option-color-hover);`)]),j(`active`,`
 color: var(--n-option-text-color-active);
 `,[i(`prefix, suffix`,`
 color: var(--n-option-text-color-active);
 `),h(`&::before`,`background-color: var(--n-option-color-active);`)]),j(`child-active`,`
 color: var(--n-option-text-color-child-active);
 `,[i(`prefix, suffix`,`
 color: var(--n-option-text-color-child-active);
 `)])]),j(`disabled`,`
 cursor: not-allowed;
 opacity: var(--n-option-opacity-disabled);
 `),j(`group`,`
 font-size: calc(var(--n-font-size) - 1px);
 color: var(--n-group-header-text-color);
 `,[i(`prefix`,`
 width: calc(var(--n-option-prefix-width) / 2);
 `,[j(`show-icon`,`
 width: calc(var(--n-option-icon-prefix-width) / 2);
 `)])]),i(`prefix`,`
 width: var(--n-option-prefix-width);
 display: flex;
 justify-content: center;
 align-items: center;
 color: var(--n-prefix-color);
 transition: color .3s var(--n-bezier);
 z-index: 1;
 `,[j(`show-icon`,`
 width: var(--n-option-icon-prefix-width);
 `),M(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),i(`label`,`
 white-space: nowrap;
 flex: 1;
 z-index: 1;
 `),i(`suffix`,`
 box-sizing: border-box;
 flex-grow: 0;
 flex-shrink: 0;
 display: flex;
 justify-content: flex-end;
 align-items: center;
 min-width: var(--n-option-suffix-width);
 padding: 0 8px;
 transition: color .3s var(--n-bezier);
 color: var(--n-suffix-color);
 z-index: 1;
 `,[j(`has-submenu`,`
 width: var(--n-option-icon-suffix-width);
 `),M(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),M(`dropdown-menu`,`pointer-events: all;`)]),M(`dropdown-offset-container`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: -4px;
 bottom: -4px;
 `)]),M(`dropdown-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 4px 0;
 `),M(`dropdown-menu-wrapper`,`
 transform-origin: var(--v-transform-origin);
 width: fit-content;
 `),h(`>`,[M(`scrollbar`,`
 height: inherit;
 max-height: inherit;
 `)]),u(`scrollable`,`
 padding: var(--n-padding);
 `),j(`scrollable`,[i(`content`,`
 padding: var(--n-padding);
 `)])]),ye={animated:{type:Boolean,default:!0},keyboard:{type:Boolean,default:!0},size:String,inverted:Boolean,placement:{type:String,default:`bottom`},onSelect:[Function,Array],options:{type:Array,default:()=>[]},menuProps:Function,showArrow:Boolean,renderLabel:Function,renderIcon:Function,renderOption:Function,nodeProps:Function,labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},value:[String,Number]},be=Object.keys(R),xe=f({name:`Dropdown`,inheritAttrs:!1,props:Object.assign(Object.assign(Object.assign({},R),ye),_.props),setup(t){let n=a(!1),i=e(r(t,`show`),n),s=d(()=>{let{keyField:e,childrenField:n}=t;return ee(t.options,{getKey(t){return t[e]},getDisabled(e){return e.disabled===!0},getIgnored(e){return e.type===`divider`||e.type===`render`},getChildren(e){return e[n]}})}),l=d(()=>s.value.treeNodes),u=a(null),f=a(null),m=a(null),h=d(()=>u.value??f.value??m.value??null),g=d(()=>s.value.getPath(h.value).keyPath),v=d(()=>s.value.getPath(t.value).keyPath),y=b(()=>t.keyboard&&i.value);W({keydown:{ArrowUp:{prevent:!0,handler:F},ArrowRight:{prevent:!0,handler:P},ArrowDown:{prevent:!0,handler:I},ArrowLeft:{prevent:!0,handler:N},Enter:{prevent:!0,handler:L},Escape:M}},y);let{mergedClsPrefixRef:x,inlineThemeDisabled:S,mergedComponentPropsRef:w}=c(t),D=d(()=>t.size||w?.value?.Dropdown?.size||`medium`),O=_(`Dropdown`,`-dropdown`,ve,q,t,x);C(Y,{labelFieldRef:r(t,`labelField`),childrenFieldRef:r(t,`childrenField`),renderLabelRef:r(t,`renderLabel`),renderIconRef:r(t,`renderIcon`),hoverKeyRef:u,keyboardKeyRef:f,lastToggledSubmenuKeyRef:m,pendingKeyPathRef:g,activeKeyPathRef:v,animatedRef:r(t,`animated`),mergedShowRef:i,nodePropsRef:r(t,`nodeProps`),renderOptionRef:r(t,`renderOption`),menuPropsRef:r(t,`menuProps`),doSelect:k,doUpdateShow:A}),T(i,e=>{!t.animated&&!e&&j()});function k(e,n){let{onSelect:r}=t;r&&E(r,e,n)}function A(e){let{"onUpdate:show":r,onUpdateShow:i}=t;r&&E(r,e),i&&E(i,e),n.value=e}function j(){u.value=null,f.value=null,m.value=null}function M(){A(!1)}function N(){z(`left`)}function P(){z(`right`)}function F(){z(`up`)}function I(){z(`down`)}function L(){let e=R();e?.isLeaf&&i.value&&(k(e.key,e.rawNode),A(!1))}function R(){let{value:e}=s,{value:t}=h;return!e||t===null?null:e.getNode(t)??null}function z(e){let{value:t}=h,{value:{getFirstAvailableNode:n}}=s,r=null;if(t===null){let e=n();e!==null&&(r=e.key)}else{let t=R();if(t){let n;switch(e){case`down`:n=t.getNext();break;case`up`:n=t.getPrev();break;case`right`:n=t.getChild();break;case`left`:n=t.getParent();break}n&&(r=n.key)}}r!==null&&(u.value=null,f.value=r)}let B=d(()=>{let{inverted:e}=t,n=D.value,{common:{cubicBezierEaseInOut:r},self:i}=O.value,{padding:a,dividerColor:s,borderRadius:c,optionOpacityDisabled:l,[o(`optionIconSuffixWidth`,n)]:u,[o(`optionSuffixWidth`,n)]:d,[o(`optionIconPrefixWidth`,n)]:f,[o(`optionPrefixWidth`,n)]:p,[o(`fontSize`,n)]:m,[o(`optionHeight`,n)]:h,[o(`optionIconSize`,n)]:g}=i,_={"--n-bezier":r,"--n-font-size":m,"--n-padding":a,"--n-border-radius":c,"--n-option-height":h,"--n-option-prefix-width":p,"--n-option-icon-prefix-width":f,"--n-option-suffix-width":d,"--n-option-icon-suffix-width":u,"--n-option-icon-size":g,"--n-divider-color":s,"--n-option-opacity-disabled":l};return e?(_[`--n-color`]=i.colorInverted,_[`--n-option-color-hover`]=i.optionColorHoverInverted,_[`--n-option-color-active`]=i.optionColorActiveInverted,_[`--n-option-text-color`]=i.optionTextColorInverted,_[`--n-option-text-color-hover`]=i.optionTextColorHoverInverted,_[`--n-option-text-color-active`]=i.optionTextColorActiveInverted,_[`--n-option-text-color-child-active`]=i.optionTextColorChildActiveInverted,_[`--n-prefix-color`]=i.prefixColorInverted,_[`--n-suffix-color`]=i.suffixColorInverted,_[`--n-group-header-text-color`]=i.groupHeaderTextColorInverted):(_[`--n-color`]=i.color,_[`--n-option-color-hover`]=i.optionColorHover,_[`--n-option-color-active`]=i.optionColorActive,_[`--n-option-text-color`]=i.optionTextColor,_[`--n-option-text-color-hover`]=i.optionTextColorHover,_[`--n-option-text-color-active`]=i.optionTextColorActive,_[`--n-option-text-color-child-active`]=i.optionTextColorChildActive,_[`--n-prefix-color`]=i.prefixColor,_[`--n-suffix-color`]=i.suffixColor,_[`--n-group-header-text-color`]=i.groupHeaderTextColor),_}),V=S?p(`dropdown`,d(()=>`${D.value[0]}${t.inverted?`i`:``}`),B,t):void 0;return{mergedClsPrefix:x,mergedTheme:O,mergedSize:D,tmNodes:l,mergedShow:i,handleAfterLeave:()=>{t.animated&&j()},doUpdateShow:A,cssVars:S?void 0:B,themeClass:V?.themeClass,onRender:V?.onRender}},render(){let e=(e,t,n,r,i)=>{var a;let{mergedClsPrefix:o,menuProps:s}=this;(a=this.onRender)==null||a.call(this);let c=s?.(void 0,this.tmNodes.map(e=>e.rawNode))||{},l={ref:G(t),class:[e,`${o}-dropdown`,`${o}-dropdown--${this.mergedSize}-size`,this.themeClass],clsPrefix:o,tmNodes:this.tmNodes,style:[...n,this.cssVars],showArrow:this.showArrow,arrowStyle:this.arrowStyle,scrollable:this.scrollable,onMouseenter:r,onMouseleave:i};return m(_e,g(this.$attrs,l,c))},{mergedTheme:t}=this,n={show:this.mergedShow,theme:t.peers.Popover,themeOverrides:t.peerOverrides.Popover,internalOnAfterLeave:this.handleAfterLeave,internalRenderBody:e,onUpdateShow:this.doUpdateShow,"onUpdate:show":void 0};return m(B,Object.assign({},oe(this.$props,be),n),{trigger:()=>{var e;return(e=this.$slots).default?.call(e)}})}});export{W as a,G as i,q as n,K as r,xe as t};