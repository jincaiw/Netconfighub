import{A as e,O as t,d as n,k as r,m as i}from"./token-DvMyJbbF.js";import{$ as a,Gt as o,Ht as s,Jn as c,Kt as l,Ln as u,Mn as d,Qt as f,Tn as p,Vt as m,Wt as h,Xn as g,Yn as _,_ as v,_t as y,bn as b,bt as x,cn as S,er as C,et as w,ft as T,g as E,gn as D,in as O,l as k,lt as A,qt as j,vn as M,wn as N,xn as P,yt as F}from"./Button-scR9_EKp.js";import{F as I,M as L,P as R,R as z,T as B,W as V,_ as H,a as U,c as ee,i as te,j as ne,l as re,q as W,s as ie,u as ae,v as oe}from"./Tooltip-DPODkf4E.js";function G(t={},n){let i=c({ctrl:!1,command:!1,win:!1,shift:!1,tab:!1}),{keydown:a,keyup:o}=t,s=e=>{switch(e.key){case`Control`:i.ctrl=!0;break;case`Meta`:i.command=!0,i.win=!0;break;case`Shift`:i.shift=!0;break;case`Tab`:i.tab=!0;break}a!==void 0&&Object.keys(a).forEach(t=>{if(t!==e.key)return;let n=a[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},l=e=>{switch(e.key){case`Control`:i.ctrl=!1;break;case`Meta`:i.command=!1,i.win=!1;break;case`Shift`:i.shift=!1;break;case`Tab`:i.tab=!1;break}o!==void 0&&Object.keys(o).forEach(t=>{if(t!==e.key)return;let n=o[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},d=()=>{(n===void 0||n.value)&&(e(`keydown`,document,s),e(`keyup`,document,l)),n!==void 0&&u(n,t=>{t?(e(`keydown`,document,s),e(`keyup`,document,l)):(r(`keydown`,document,s),r(`keyup`,document,l))})};return V()?(N(d),p(()=>{(n===void 0||n.value)&&(r(`keydown`,document,s),r(`keyup`,document,l))})):d(),_(i)}function se(e,t,n){if(!t)return e;let r=g(e.value),i=null;return u(e,e=>{i!==null&&window.clearTimeout(i),e===!0?n&&!n.value?r.value=!0:i=window.setTimeout(()=>{r.value=!0},t):r.value=!1}),r}function K(e){return t=>{t?e.value=t.$el:e.value=null}}var q=D({name:`ChevronRight`,render(){return M(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},M(`path`,{d:`M5.64645 3.14645C5.45118 3.34171 5.45118 3.65829 5.64645 3.85355L9.79289 8L5.64645 12.1464C5.45118 12.3417 5.45118 12.6583 5.64645 12.8536C5.84171 13.0488 6.15829 13.0488 6.35355 12.8536L10.8536 8.35355C11.0488 8.15829 11.0488 7.84171 10.8536 7.64645L6.35355 3.14645C6.15829 2.95118 5.84171 2.95118 5.64645 3.14645Z`,fill:`currentColor`}))}}),ce={padding:`4px 0`,optionIconSizeSmall:`14px`,optionIconSizeMedium:`16px`,optionIconSizeLarge:`16px`,optionIconSizeHuge:`18px`,optionSuffixWidthSmall:`14px`,optionSuffixWidthMedium:`14px`,optionSuffixWidthLarge:`16px`,optionSuffixWidthHuge:`16px`,optionIconSuffixWidthSmall:`32px`,optionIconSuffixWidthMedium:`32px`,optionIconSuffixWidthLarge:`36px`,optionIconSuffixWidthHuge:`36px`,optionPrefixWidthSmall:`14px`,optionPrefixWidthMedium:`14px`,optionPrefixWidthLarge:`16px`,optionPrefixWidthHuge:`16px`,optionIconPrefixWidthSmall:`36px`,optionIconPrefixWidthMedium:`36px`,optionIconPrefixWidthLarge:`40px`,optionIconPrefixWidthHuge:`40px`};function le(e){let{primaryColor:t,textColor2:n,dividerColor:r,hoverColor:i,popoverColor:a,invertedColor:o,borderRadius:s,fontSizeSmall:c,fontSizeMedium:l,fontSizeLarge:u,fontSizeHuge:d,heightSmall:f,heightMedium:p,heightLarge:m,heightHuge:h,textColor3:g,opacityDisabled:_}=e;return Object.assign(Object.assign({},ce),{optionHeightSmall:f,optionHeightMedium:p,optionHeightLarge:m,optionHeightHuge:h,borderRadius:s,fontSizeSmall:c,fontSizeMedium:l,fontSizeLarge:u,fontSizeHuge:d,optionTextColor:n,optionTextColorHover:n,optionTextColorActive:t,optionTextColorChildActive:t,color:a,dividerColor:r,suffixColor:n,prefixColor:n,optionColorHover:i,optionColorActive:x(t,{alpha:.1}),groupHeaderTextColor:g,optionTextColorInverted:`#BBB`,optionTextColorHoverInverted:`#FFF`,optionTextColorActiveInverted:`#FFF`,optionTextColorChildActiveInverted:`#FFF`,colorInverted:o,dividerColorInverted:`#BBB`,suffixColorInverted:`#BBB`,prefixColorInverted:`#BBB`,optionColorHoverInverted:t,optionColorActiveInverted:t,groupHeaderTextColorInverted:`#AAA`,optionOpacityDisabled:_})}var J=E({name:`Dropdown`,common:k,peers:{Popover:ee},self:le}),Y=y(`n-dropdown-menu`),X=y(`n-dropdown`),Z=y(`n-dropdown-option`),Q=D({name:`DropdownDivider`,props:{clsPrefix:{type:String,required:!0}},render(){return M(`div`,{class:`${this.clsPrefix}-dropdown-divider`})}}),ue=D({name:`DropdownGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{showIconRef:e,hasSubmenuRef:t}=b(Y),{renderLabelRef:n,labelFieldRef:r,nodePropsRef:i,renderOptionRef:a}=b(X);return{labelField:r,showIcon:e,hasSubmenu:t,renderLabel:n,nodeProps:i,renderOption:a}},render(){let{clsPrefix:e,hasSubmenu:t,showIcon:n,nodeProps:r,renderLabel:i,renderOption:a}=this,{rawNode:o}=this.tmNode,s=M(`div`,Object.assign({class:`${e}-dropdown-option`},r?.(o)),M(`div`,{class:`${e}-dropdown-option-body ${e}-dropdown-option-body--group`},M(`div`,{"data-dropdown-option":!0,class:[`${e}-dropdown-option-body__prefix`,n&&`${e}-dropdown-option-body__prefix--show-icon`]},H(o.icon)),M(`div`,{class:`${e}-dropdown-option-body__label`,"data-dropdown-option":!0},i?i(o):H(o.title??o[this.labelField])),M(`div`,{class:[`${e}-dropdown-option-body__suffix`,t&&`${e}-dropdown-option-body__suffix--has-submenu`],"data-dropdown-option":!0})));return a?a({node:s,option:o}):s}});function $(e,t){return e.type===`submenu`||e.type===void 0&&e[t]!==void 0}function de(e){return e.type===`group`}function fe(e){return e.type===`divider`}function pe(e){return e.type===`render`}var me=D({name:`DropdownOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null},placement:{type:String,default:`right-start`},props:Object,scrollable:Boolean},setup(e){let t=b(X),{hoverKeyRef:n,keyboardKeyRef:r,lastToggledSubmenuKeyRef:i,pendingKeyPathRef:a,activeKeyPathRef:o,animatedRef:s,mergedShowRef:c,renderLabelRef:l,renderIconRef:u,labelFieldRef:f,childrenFieldRef:p,renderOptionRef:m,nodePropsRef:h,menuPropsRef:_}=t,v=b(Z,null),y=b(Y),x=b(R),C=S(()=>e.tmNode.rawNode),w=S(()=>{let{value:t}=p;return $(e.tmNode.rawNode,t)}),T=S(()=>{let{disabled:t}=e.tmNode;return t}),E=se(S(()=>{if(!w.value)return!1;let{key:t,disabled:o}=e.tmNode;if(o)return!1;let{value:s}=n,{value:c}=r,{value:l}=i,{value:u}=a;return s===null?c===null?l===null?!1:u.includes(t):u.includes(t)&&u[u.length-1]!==t:u.includes(t)}),300,S(()=>r.value===null&&!s.value)),D=S(()=>!!v?.enteringSubmenuRef.value),O=g(!1);d(Z,{enteringSubmenuRef:O});function k(){O.value=!0}function A(){O.value=!1}function j(){let{parentKey:t,tmNode:a}=e;a.disabled||c.value&&(i.value=t,r.value=null,n.value=a.key)}function M(){let{tmNode:t}=e;t.disabled||c.value&&n.value!==t.key&&j()}function N(t){if(e.tmNode.disabled||!c.value)return;let{relatedTarget:r}=t;r&&!W({target:r},`dropdownOption`)&&!W({target:r},`scrollbarRail`)&&(n.value=null)}function P(){let{value:n}=w,{tmNode:r}=e;c.value&&!n&&!r.disabled&&(t.doSelect(r.key,r.rawNode),t.doUpdateShow(!1))}return{labelField:f,renderLabel:l,renderIcon:u,siblingHasIcon:y.showIconRef,siblingHasSubmenu:y.hasSubmenuRef,menuProps:_,popoverBody:x,animated:s,mergedShowSubmenu:S(()=>E.value&&!D.value),rawNode:C,hasSubmenu:w,pending:F(()=>{let{value:t}=a,{key:n}=e.tmNode;return t.includes(n)}),childActive:F(()=>{let{value:t}=o,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r===-1?!1:r<t.length-1}),active:F(()=>{let{value:t}=o,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r===-1?!1:r===t.length-1}),mergedDisabled:T,renderOption:m,nodeProps:h,handleClick:P,handleMouseMove:M,handleMouseEnter:j,handleMouseLeave:N,handleSubmenuBeforeEnter:k,handleSubmenuAfterEnter:A}},render(){let{animated:e,rawNode:t,mergedShowSubmenu:r,clsPrefix:i,siblingHasIcon:a,siblingHasSubmenu:o,renderLabel:s,renderIcon:c,renderOption:l,nodeProps:u,props:d,scrollable:p}=this,m=null;if(r){let e=this.menuProps?.call(this,t,t.children);m=M(_e,Object.assign({},e,{clsPrefix:i,scrollable:this.scrollable,tmNodes:this.tmNode.children,parentKey:this.tmNode.key}))}let h={class:[`${i}-dropdown-option-body`,this.pending&&`${i}-dropdown-option-body--pending`,this.active&&`${i}-dropdown-option-body--active`,this.childActive&&`${i}-dropdown-option-body--child-active`,this.mergedDisabled&&`${i}-dropdown-option-body--disabled`],onMousemove:this.handleMouseMove,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onClick:this.handleClick},g=u?.(t),_=M(`div`,Object.assign({class:[`${i}-dropdown-option`,g?.class],"data-dropdown-option":!0},g),M(`div`,P(h,d),[M(`div`,{class:[`${i}-dropdown-option-body__prefix`,a&&`${i}-dropdown-option-body__prefix--show-icon`]},[c?c(t):H(t.icon)]),M(`div`,{"data-dropdown-option":!0,class:`${i}-dropdown-option-body__label`},s?s(t):H(t[this.labelField]??t.title)),M(`div`,{"data-dropdown-option":!0,class:[`${i}-dropdown-option-body__suffix`,o&&`${i}-dropdown-option-body__suffix--has-submenu`]},this.hasSubmenu?M(n,null,{default:()=>M(q,null)}):null)]),this.hasSubmenu?M(L,null,{default:()=>[M(ne,null,{default:()=>M(`div`,{class:`${i}-dropdown-offset-container`},M(B,{show:this.mergedShowSubmenu,placement:this.placement,to:p&&this.popoverBody||void 0,teleportDisabled:!p},{default:()=>M(`div`,{class:`${i}-dropdown-menu-wrapper`},e?M(f,{onBeforeEnter:this.handleSubmenuBeforeEnter,onAfterEnter:this.handleSubmenuAfterEnter,name:`fade-in-scale-up-transition`,appear:!0},{default:()=>m}):m)}))})]}):null);return l?l({node:_,option:t}):_}}),he=D({name:`NDropdownGroup`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null}},render(){let{tmNode:e,parentKey:t,clsPrefix:n}=this,{children:r}=e;return M(O,null,M(ue,{clsPrefix:n,tmNode:e,key:e.key}),r?.map(e=>{let{rawNode:r}=e;return r.show===!1?null:fe(r)?M(Q,{clsPrefix:n,key:e.key}):e.isGroup?(T(`dropdown`,"`group` node is not allowed to be put in `group` node."),null):M(me,{clsPrefix:n,tmNode:e,parentKey:t,key:e.key})}))}}),ge=D({name:`DropdownRenderOption`,props:{tmNode:{type:Object,required:!0}},render(){let{rawNode:{render:e,props:t}}=this.tmNode;return M(`div`,t,[e?.()])}}),_e=D({name:`DropdownMenu`,props:{scrollable:Boolean,showArrow:Boolean,arrowStyle:[String,Object],clsPrefix:{type:String,required:!0},tmNodes:{type:Array,default:()=>[]},parentKey:{type:[String,Number],default:null}},setup(e){let{renderIconRef:t,childrenFieldRef:n}=b(X);d(Y,{showIconRef:S(()=>{let n=t.value;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>n?n(e):e.icon);let{rawNode:t}=e;return n?n(t):t.icon})}),hasSubmenuRef:S(()=>{let{value:t}=n;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>$(e,t));let{rawNode:n}=e;return $(n,t)})})});let r=g(null);return d(I,null),d(z,null),d(R,r),{bodyRef:r}},render(){let{parentKey:e,clsPrefix:t,scrollable:n}=this,r=this.tmNodes.map(r=>{let{rawNode:i}=r;return i.show===!1?null:pe(i)?M(ge,{tmNode:r,key:r.key}):fe(i)?M(Q,{clsPrefix:t,key:r.key}):de(i)?M(he,{clsPrefix:t,tmNode:r,parentKey:e,key:r.key}):M(me,{clsPrefix:t,tmNode:r,parentKey:e,key:r.key,props:i.props,scrollable:n})});return M(`div`,{class:[`${t}-dropdown-menu`,n&&`${t}-dropdown-menu--scrollable`],ref:`bodyRef`},n?M(i,{contentClass:`${t}-dropdown-menu__content`},{default:()=>r}):r,this.showArrow?ie({clsPrefix:t,arrowStyle:this.arrowStyle,arrowClass:void 0,arrowWrapperClass:void 0,arrowWrapperStyle:void 0}):null)}}),ve=s(`dropdown-menu`,`
 transform-origin: var(--v-transform-origin);
 background-color: var(--n-color);
 border-radius: var(--n-border-radius);
 box-shadow: var(--n-box-shadow);
 position: relative;
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
`,[re(),s(`dropdown-option`,`
 position: relative;
 `,[m(`a`,`
 text-decoration: none;
 color: inherit;
 outline: none;
 `,[m(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),s(`dropdown-option-body`,`
 display: flex;
 cursor: pointer;
 position: relative;
 height: var(--n-option-height);
 line-height: var(--n-option-height);
 font-size: var(--n-font-size);
 color: var(--n-option-text-color);
 transition: color .3s var(--n-bezier);
 `,[m(`&::before`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 left: 4px;
 right: 4px;
 transition: background-color .3s var(--n-bezier);
 border-radius: var(--n-border-radius);
 `),l(`disabled`,[o(`pending`,`
 color: var(--n-option-text-color-hover);
 `,[h(`prefix, suffix`,`
 color: var(--n-option-text-color-hover);
 `),m(`&::before`,`background-color: var(--n-option-color-hover);`)]),o(`active`,`
 color: var(--n-option-text-color-active);
 `,[h(`prefix, suffix`,`
 color: var(--n-option-text-color-active);
 `),m(`&::before`,`background-color: var(--n-option-color-active);`)]),o(`child-active`,`
 color: var(--n-option-text-color-child-active);
 `,[h(`prefix, suffix`,`
 color: var(--n-option-text-color-child-active);
 `)])]),o(`disabled`,`
 cursor: not-allowed;
 opacity: var(--n-option-opacity-disabled);
 `),o(`group`,`
 font-size: calc(var(--n-font-size) - 1px);
 color: var(--n-group-header-text-color);
 `,[h(`prefix`,`
 width: calc(var(--n-option-prefix-width) / 2);
 `,[o(`show-icon`,`
 width: calc(var(--n-option-icon-prefix-width) / 2);
 `)])]),h(`prefix`,`
 width: var(--n-option-prefix-width);
 display: flex;
 justify-content: center;
 align-items: center;
 color: var(--n-prefix-color);
 transition: color .3s var(--n-bezier);
 z-index: 1;
 `,[o(`show-icon`,`
 width: var(--n-option-icon-prefix-width);
 `),s(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),h(`label`,`
 white-space: nowrap;
 flex: 1;
 z-index: 1;
 `),h(`suffix`,`
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
 `,[o(`has-submenu`,`
 width: var(--n-option-icon-suffix-width);
 `),s(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),s(`dropdown-menu`,`pointer-events: all;`)]),s(`dropdown-offset-container`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: -4px;
 bottom: -4px;
 `)]),s(`dropdown-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 4px 0;
 `),s(`dropdown-menu-wrapper`,`
 transform-origin: var(--v-transform-origin);
 width: fit-content;
 `),m(`>`,[s(`scrollbar`,`
 height: inherit;
 max-height: inherit;
 `)]),l(`scrollable`,`
 padding: var(--n-padding);
 `),o(`scrollable`,[h(`content`,`
 padding: var(--n-padding);
 `)])]),ye={animated:{type:Boolean,default:!0},keyboard:{type:Boolean,default:!0},size:String,inverted:Boolean,placement:{type:String,default:`bottom`},onSelect:[Function,Array],options:{type:Array,default:()=>[]},menuProps:Function,showArrow:Boolean,renderLabel:Function,renderIcon:Function,renderOption:Function,nodeProps:Function,labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},value:[String,Number]},be=Object.keys(U),xe=Object.assign(Object.assign(Object.assign({},U),ye),v.props),Se=D({name:`Dropdown`,inheritAttrs:!1,props:xe,setup(e){let n=g(!1),r=t(C(e,`show`),n),i=S(()=>{let{keyField:t,childrenField:n}=e;return ae(e.options,{getKey(e){return e[t]},getDisabled(e){return e.disabled===!0},getIgnored(e){return e.type===`divider`||e.type===`render`},getChildren(e){return e[n]}})}),o=S(()=>i.value.treeNodes),s=g(null),c=g(null),l=g(null),f=S(()=>s.value??c.value??l.value??null),p=S(()=>i.value.getPath(f.value).keyPath),m=S(()=>i.value.getPath(e.value).keyPath),h=F(()=>e.keyboard&&r.value);G({keydown:{ArrowUp:{prevent:!0,handler:P},ArrowRight:{prevent:!0,handler:N},ArrowDown:{prevent:!0,handler:I},ArrowLeft:{prevent:!0,handler:M},Enter:{prevent:!0,handler:L},Escape:k}},h);let{mergedClsPrefixRef:_,inlineThemeDisabled:y,mergedComponentPropsRef:b}=w(e),x=S(()=>e.size||b?.value?.Dropdown?.size||`medium`),T=v(`Dropdown`,`-dropdown`,ve,J,e,_);d(X,{labelFieldRef:C(e,`labelField`),childrenFieldRef:C(e,`childrenField`),renderLabelRef:C(e,`renderLabel`),renderIconRef:C(e,`renderIcon`),hoverKeyRef:s,keyboardKeyRef:c,lastToggledSubmenuKeyRef:l,pendingKeyPathRef:p,activeKeyPathRef:m,animatedRef:C(e,`animated`),mergedShowRef:r,nodePropsRef:C(e,`nodeProps`),renderOptionRef:C(e,`renderOption`),menuPropsRef:C(e,`menuProps`),doSelect:E,doUpdateShow:D}),u(r,t=>{!e.animated&&!t&&O()});function E(t,n){let{onSelect:r}=e;r&&A(r,t,n)}function D(t){let{"onUpdate:show":r,onUpdateShow:i}=e;r&&A(r,t),i&&A(i,t),n.value=t}function O(){s.value=null,c.value=null,l.value=null}function k(){D(!1)}function M(){z(`left`)}function N(){z(`right`)}function P(){z(`up`)}function I(){z(`down`)}function L(){let e=R();e?.isLeaf&&r.value&&(E(e.key,e.rawNode),D(!1))}function R(){let{value:e}=i,{value:t}=f;return!e||t===null?null:e.getNode(t)??null}function z(e){let{value:t}=f,{value:{getFirstAvailableNode:n}}=i,r=null;if(t===null){let e=n();e!==null&&(r=e.key)}else{let t=R();if(t){let n;switch(e){case`down`:n=t.getNext();break;case`up`:n=t.getPrev();break;case`right`:n=t.getChild();break;case`left`:n=t.getParent();break}n&&(r=n.key)}}r!==null&&(s.value=null,c.value=r)}let B=S(()=>{let{inverted:t}=e,n=x.value,{common:{cubicBezierEaseInOut:r},self:i}=T.value,{padding:a,dividerColor:o,borderRadius:s,optionOpacityDisabled:c,[j(`optionIconSuffixWidth`,n)]:l,[j(`optionSuffixWidth`,n)]:u,[j(`optionIconPrefixWidth`,n)]:d,[j(`optionPrefixWidth`,n)]:f,[j(`fontSize`,n)]:p,[j(`optionHeight`,n)]:m,[j(`optionIconSize`,n)]:h}=i,g={"--n-bezier":r,"--n-font-size":p,"--n-padding":a,"--n-border-radius":s,"--n-option-height":m,"--n-option-prefix-width":f,"--n-option-icon-prefix-width":d,"--n-option-suffix-width":u,"--n-option-icon-suffix-width":l,"--n-option-icon-size":h,"--n-divider-color":o,"--n-option-opacity-disabled":c};return t?(g[`--n-color`]=i.colorInverted,g[`--n-option-color-hover`]=i.optionColorHoverInverted,g[`--n-option-color-active`]=i.optionColorActiveInverted,g[`--n-option-text-color`]=i.optionTextColorInverted,g[`--n-option-text-color-hover`]=i.optionTextColorHoverInverted,g[`--n-option-text-color-active`]=i.optionTextColorActiveInverted,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActiveInverted,g[`--n-prefix-color`]=i.prefixColorInverted,g[`--n-suffix-color`]=i.suffixColorInverted,g[`--n-group-header-text-color`]=i.groupHeaderTextColorInverted):(g[`--n-color`]=i.color,g[`--n-option-color-hover`]=i.optionColorHover,g[`--n-option-color-active`]=i.optionColorActive,g[`--n-option-text-color`]=i.optionTextColor,g[`--n-option-text-color-hover`]=i.optionTextColorHover,g[`--n-option-text-color-active`]=i.optionTextColorActive,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActive,g[`--n-prefix-color`]=i.prefixColor,g[`--n-suffix-color`]=i.suffixColor,g[`--n-group-header-text-color`]=i.groupHeaderTextColor),g}),V=y?a(`dropdown`,S(()=>`${x.value[0]}${e.inverted?`i`:``}`),B,e):void 0;return{mergedClsPrefix:_,mergedTheme:T,mergedSize:x,tmNodes:o,mergedShow:r,handleAfterLeave:()=>{e.animated&&O()},doUpdateShow:D,cssVars:y?void 0:B,themeClass:V?.themeClass,onRender:V?.onRender}},render(){let e=(e,t,n,r,i)=>{var a;let{mergedClsPrefix:o,menuProps:s}=this;(a=this.onRender)==null||a.call(this);let c=s?.(void 0,this.tmNodes.map(e=>e.rawNode))||{},l={ref:K(t),class:[e,`${o}-dropdown`,`${o}-dropdown--${this.mergedSize}-size`,this.themeClass],clsPrefix:o,tmNodes:this.tmNodes,style:[...n,this.cssVars],showArrow:this.showArrow,arrowStyle:this.arrowStyle,scrollable:this.scrollable,onMouseenter:r,onMouseleave:i};return M(_e,P(this.$attrs,l,c))},{mergedTheme:t}=this,n={show:this.mergedShow,theme:t.peers.Popover,themeOverrides:t.peerOverrides.Popover,internalOnAfterLeave:this.handleAfterLeave,internalRenderBody:e,onUpdateShow:this.doUpdateShow,"onUpdate:show":void 0};return M(te,Object.assign({},oe(this.$props,be),n),{trigger:()=>{var e;return(e=this.$slots).default?.call(e)}})}});export{K as a,q as i,xe as n,G as o,J as r,Se as t};