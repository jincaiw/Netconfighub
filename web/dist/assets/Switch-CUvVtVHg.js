import{m as e}from"./api-CZuskCjY.js";import{An as t,Ct as n,En as r,Et as i,Q as a,Tt as o,Vt as s,X as c,Yt as l,Z as u,Zt as d,at as f,bt as p,c as m,d as h,g,gt as _,h as v,ht as y,m as b,mn as x,n as S,nt as C,ot as w,p as T,r as E,rt as D,s as O,tn as ee,vt as te,wt as k,xt as A,y as j}from"./Button-v7qKim5p.js";import{D as M,w as N,x as ne}from"./token-dQuugHoF.js";import{n as P,o as re,t as F}from"./Input-C6Yqttcz.js";import{t as I}from"./Add-DnwwuZfa.js";var ie=l({name:`Remove`,render(){return d(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 512 512`},d(`line`,{x1:`400`,y1:`256`,x2:`112`,y2:`256`,style:`
        fill: none;
        stroke: currentColor;
        stroke-linecap: round;
        stroke-linejoin: round;
        stroke-width: 32px;
      `}))}});function L(e){let{textColorDisabled:t}=e;return{iconColorDisabled:t}}var ae=v({name:`InputNumber`,common:O,peers:{Button:E,Input:P},self:L}),R={buttonHeightSmall:`14px`,buttonHeightMedium:`18px`,buttonHeightLarge:`22px`,buttonWidthSmall:`14px`,buttonWidthMedium:`18px`,buttonWidthLarge:`22px`,buttonWidthPressedSmall:`20px`,buttonWidthPressedMedium:`24px`,buttonWidthPressedLarge:`28px`,railHeightSmall:`18px`,railHeightMedium:`22px`,railHeightLarge:`26px`,railWidthSmall:`32px`,railWidthMedium:`40px`,railWidthLarge:`48px`};function z(e){let{primaryColor:t,opacityDisabled:n,borderRadius:r,textColor3:i}=e;return Object.assign(Object.assign({},R),{iconColor:i,textColor:`white`,loadingColor:t,opacityDisabled:n,railColor:`rgba(0, 0, 0, .14)`,railColorActive:t,buttonBoxShadow:`0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)`,buttonColor:`#FFF`,railBorderRadiusSmall:r,railBorderRadiusMedium:r,railBorderRadiusLarge:r,buttonBorderRadiusSmall:r,buttonBorderRadiusMedium:r,buttonBorderRadiusLarge:r,boxShadowFocus:`0 0 0 2px ${_(t,{alpha:.2})}`})}var oe={name:`Switch`,common:O,self:z},se=p([A(`input-number-suffix`,`
 display: inline-block;
 margin-right: 10px;
 `),A(`input-number-prefix`,`
 display: inline-block;
 margin-left: 10px;
 `)]);function ce(e){return e==null||typeof e==`string`&&e.trim()===``?null:Number(e)}function le(e){return e.includes(`.`)&&(/^(-)?\d+.*(\.|0)$/.test(e)||/^-?\d*$/.test(e))||e===`-`||e===`-0`}function B(e){return e==null?!0:!Number.isNaN(e)}function ue(e,t){return typeof e==`number`?t===void 0?String(e):e.toFixed(t):``}function V(e){if(e===null)return null;if(typeof e==`number`)return e;{let t=Number(e);return Number.isNaN(t)?null:t}}var H=800,U=100,de=l({name:`InputNumber`,props:Object.assign(Object.assign({},g.props),{autofocus:Boolean,loading:{type:Boolean,default:void 0},placeholder:String,defaultValue:{type:Number,default:null},value:Number,step:{type:[Number,String],default:1},min:[Number,String],max:[Number,String],size:String,disabled:{type:Boolean,default:void 0},validator:Function,bordered:{type:Boolean,default:void 0},showButton:{type:Boolean,default:!0},buttonPlacement:{type:String,default:`right`},inputProps:Object,readonly:Boolean,clearable:Boolean,keyboard:{type:Object,default:{}},updateValueOnInput:{type:Boolean,default:!0},round:{type:Boolean,default:void 0},parse:Function,format:Function,precision:Number,status:String,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onFocus:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onChange:[Function,Array]}),slots:Object,setup(n){let{mergedBorderedRef:i,mergedClsPrefixRef:o,mergedRtlRef:l,mergedComponentPropsRef:u}=a(n),d=g(`InputNumber`,`-input-number`,se,ae,n,o),{localeRef:f}=re(`InputNumber`),p=c(n,{mergedSize:e=>{let{size:t}=n;if(t)return t;let{mergedSize:r}=e||{};return r?.value?r.value:u?.value?.InputNumber?.size||`medium`}}),{mergedSizeRef:m,mergedDisabledRef:h,mergedStatusRef:_}=p,v=r(null),b=r(null),S=r(null),C=r(n.defaultValue),T=e(t(n,`value`),C),E=r(``),D=e=>{let t=String(e).split(`.`)[1];return t?t.length:0},O=e=>{let t=[n.min,n.max,n.step,e].map(e=>e===void 0?0:D(e));return Math.max(...t)},k=y(()=>{let{placeholder:e}=n;return e===void 0?f.value.placeholder:e}),A=y(()=>{let e=V(n.step);return e===null||e===0?1:Math.abs(e)}),M=y(()=>{let e=V(n.min);return e===null?null:e}),N=y(()=>{let e=V(n.max);return e===null?null:e}),P=()=>{let{value:e}=T;if(B(e)){let{format:t,precision:r}=n;t?E.value=t(e):e===null||r===void 0||D(e)>r?E.value=ue(e,void 0):E.value=ue(e,r)}else E.value=String(e)};P();let F=e=>{let{value:t}=T;if(e===t){P();return}let{"onUpdate:value":r,onUpdateValue:i,onChange:a}=n,{nTriggerFormInput:o,nTriggerFormChange:s}=p;a&&w(a,e),i&&w(i,e),r&&w(r,e),C.value=e,o(),s()},I=({offset:e,doUpdateIfValid:t,fixPrecision:r,isInputing:i})=>{let{value:a}=E;if(i&&le(a))return!1;let o=(n.parse||ce)(a);if(o===null)return t&&F(null),null;if(B(o)){let a=D(o),{precision:s}=n;if(s!==void 0&&s<a&&!r)return!1;let c=Number.parseFloat((o+e).toFixed(s??O(o)));if(B(c)){let{value:e}=N,{value:r}=M;if(e!==null&&c>e){if(!t||i)return!1;c=e}if(r!==null&&c<r){if(!t||i)return!1;c=r}return n.validator&&!n.validator(c)?!1:(t&&F(c),c)}}return!1},ie=y(()=>I({offset:0,doUpdateIfValid:!1,isInputing:!1,fixPrecision:!1})===!1),L=y(()=>{let{value:e}=T;if(n.validator&&e===null)return!1;let{value:t}=A;return I({offset:-t,doUpdateIfValid:!1,isInputing:!1,fixPrecision:!1})!==!1}),R=y(()=>{let{value:e}=T;if(n.validator&&e===null)return!1;let{value:t}=A;return I({offset:+t,doUpdateIfValid:!1,isInputing:!1,fixPrecision:!1})!==!1});function z(e){let{onFocus:t}=n,{nTriggerFormFocus:r}=p;t&&w(t,e),r()}function oe(e){if(e.target===v.value?.wrapperElRef)return;let t=I({offset:0,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0});if(t!==!1){let e=v.value?.inputElRef;e&&(e.value=String(t||``)),T.value===t&&P()}else P();let{onBlur:r}=n,{nTriggerFormBlur:i}=p;r&&w(r,e),i(),ee(()=>{P()})}function de(e){let{onClear:t}=n;t&&w(t,e)}function W(){let{value:e}=R;if(!e){$();return}let{value:t}=T;if(t===null)n.validator||F(q());else{let{value:e}=A;I({offset:e,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0})}}function G(){let{value:e}=L;if(!e){Z();return}let{value:t}=T;if(t===null)n.validator||F(q());else{let{value:e}=A;I({offset:-e,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0})}}let K=z,fe=oe;function q(){if(n.validator)return null;let{value:e}=M,{value:t}=N;return e===null?t===null?0:Math.min(0,t):Math.max(0,e)}function pe(e){de(e),F(null)}function me(e){var t;S.value?.$el.contains(e.target)&&e.preventDefault(),b.value?.$el.contains(e.target)&&e.preventDefault(),(t=v.value)==null||t.activate()}let J=null,Y=null,X=null;function Z(){X&&=(window.clearTimeout(X),null),J&&=(window.clearInterval(J),null)}let Q=null;function $(){Q&&=(window.clearTimeout(Q),null),Y&&=(window.clearInterval(Y),null)}function he(){Z(),X=window.setTimeout(()=>{J=window.setInterval(()=>{G()},U)},H),ne(`mouseup`,document,Z,{once:!0})}function ge(){$(),Q=window.setTimeout(()=>{Y=window.setInterval(()=>{W()},U)},H),ne(`mouseup`,document,$,{once:!0})}let _e=()=>{Y||W()},ve=()=>{J||G()};function ye(e){var t;if(e.key===`Enter`){if(e.target===v.value?.wrapperElRef)return;I({offset:0,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0})!==!1&&((t=v.value)==null||t.deactivate())}else if(e.key===`ArrowUp`){if(!R.value||n.keyboard.ArrowUp===!1)return;e.preventDefault(),I({offset:0,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0})!==!1&&W()}else if(e.key===`ArrowDown`){if(!L.value||n.keyboard.ArrowDown===!1)return;e.preventDefault(),I({offset:0,doUpdateIfValid:!0,isInputing:!1,fixPrecision:!0})!==!1&&G()}}function be(e){E.value=e,n.updateValueOnInput&&!n.format&&!n.parse&&n.precision===void 0&&I({offset:0,doUpdateIfValid:!0,isInputing:!0,fixPrecision:!1})}x(T,()=>{P()});let xe={focus:()=>v.value?.focus(),blur:()=>v.value?.blur(),select:()=>v.value?.select()},Se=j(`InputNumber`,l,o);return Object.assign(Object.assign({},xe),{rtlEnabled:Se,inputInstRef:v,minusButtonInstRef:b,addButtonInstRef:S,mergedClsPrefix:o,mergedBordered:i,uncontrolledValue:C,mergedValue:T,mergedPlaceholder:k,displayedValueInvalid:ie,mergedSize:m,mergedDisabled:h,displayedValue:E,addable:R,minusable:L,mergedStatus:_,handleFocus:K,handleBlur:fe,handleClear:pe,handleMouseDown:me,handleAddClick:_e,handleMinusClick:ve,handleAddMousedown:ge,handleMinusMousedown:he,handleKeyDown:ye,handleUpdateDisplayedValue:be,mergedTheme:d,inputThemeOverrides:{paddingSmall:`0 8px 0 10px`,paddingMedium:`0 8px 0 12px`,paddingLarge:`0 8px 0 14px`},buttonThemeOverrides:s(()=>{let{self:{iconColorDisabled:e}}=d.value,[t,n,r,i]=te(e);return{textColorTextDisabled:`rgb(${t}, ${n}, ${r})`,opacityDisabled:`${i}`}})})},render(){let{mergedClsPrefix:e,$slots:t}=this,n=()=>d(S,{text:!0,disabled:!this.minusable||this.mergedDisabled||this.readonly,focusable:!1,theme:this.mergedTheme.peers.Button,themeOverrides:this.mergedTheme.peerOverrides.Button,builtinThemeOverrides:this.buttonThemeOverrides,onClick:this.handleMinusClick,onMousedown:this.handleMinusMousedown,ref:`minusButtonInstRef`},{icon:()=>D(t[`minus-icon`],()=>[d(b,{clsPrefix:e},{default:()=>d(ie,null)})])}),r=()=>d(S,{text:!0,disabled:!this.addable||this.mergedDisabled||this.readonly,focusable:!1,theme:this.mergedTheme.peers.Button,themeOverrides:this.mergedTheme.peerOverrides.Button,builtinThemeOverrides:this.buttonThemeOverrides,onClick:this.handleAddClick,onMousedown:this.handleAddMousedown,ref:`addButtonInstRef`},{icon:()=>D(t[`add-icon`],()=>[d(b,{clsPrefix:e},{default:()=>d(I,null)})])});return d(`div`,{class:[`${e}-input-number`,this.rtlEnabled&&`${e}-input-number--rtl`]},d(F,{ref:`inputInstRef`,autofocus:this.autofocus,status:this.mergedStatus,bordered:this.mergedBordered,loading:this.loading,value:this.displayedValue,onUpdateValue:this.handleUpdateDisplayedValue,theme:this.mergedTheme.peers.Input,themeOverrides:this.mergedTheme.peerOverrides.Input,builtinThemeOverrides:this.inputThemeOverrides,size:this.mergedSize,placeholder:this.mergedPlaceholder,disabled:this.mergedDisabled,readonly:this.readonly,round:this.round,textDecoration:this.displayedValueInvalid?`line-through`:void 0,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onClear:this.handleClear,clearable:this.clearable,inputProps:this.inputProps,internalLoadingBeforeSuffix:!0},{prefix:()=>this.showButton&&this.buttonPlacement===`both`?[n(),f(t.prefix,t=>t?d(`span`,{class:`${e}-input-number-prefix`},t):null)]:t.prefix?.call(t),suffix:()=>this.showButton?[f(t.suffix,t=>t?d(`span`,{class:`${e}-input-number-suffix`},t):null),this.buttonPlacement===`right`?n():null,r()]:t.suffix?.call(t)}))}}),W=A(`switch`,`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[n(`children-placeholder`,`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),n(`rail-placeholder`,`
 display: flex;
 flex-wrap: none;
 `),n(`button-placeholder`,`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),A(`base-loading`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[h({left:`50%`,top:`50%`,originalTransform:`translateX(-50%) translateY(-50%)`})]),n(`checked, unchecked`,`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),n(`checked`,`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),n(`unchecked`,`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),p(`&:focus`,[n(`rail`,`
 box-shadow: var(--n-box-shadow-focus);
 `)]),k(`round`,[n(`rail`,`border-radius: calc(var(--n-rail-height) / 2);`,[n(`button`,`border-radius: calc(var(--n-button-height) / 2);`)])]),o(`disabled`,[o(`icon`,[k(`rubber-band`,[k(`pressed`,[n(`rail`,[n(`button`,`max-width: var(--n-button-width-pressed);`)])]),n(`rail`,[p(`&:active`,[n(`button`,`max-width: var(--n-button-width-pressed);`)])]),k(`active`,[k(`pressed`,[n(`rail`,[n(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])]),n(`rail`,[p(`&:active`,[n(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])])])])])]),k(`active`,[n(`rail`,[n(`button`,`left: calc(100% - var(--n-button-width) - var(--n-offset))`)])]),n(`rail`,`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[n(`button-icon`,`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[h()]),n(`button`,`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),k(`active`,[n(`rail`,`background-color: var(--n-rail-color-active);`)]),k(`loading`,[n(`rail`,`
 cursor: wait;
 `)]),k(`disabled`,[n(`rail`,`
 cursor: not-allowed;
 opacity: .5;
 `)])]),G=Object.assign(Object.assign({},g.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]}),K,fe=l({name:`Switch`,props:G,slots:Object,setup(n){K===void 0&&(K=typeof CSS<`u`?CSS.supports===void 0?!1:CSS.supports(`width`,`max(1px)`):!0);let{mergedClsPrefixRef:o,inlineThemeDisabled:l,mergedComponentPropsRef:d}=a(n),f=g(`Switch`,`-switch`,W,oe,n,o),p=c(n,{mergedSize(e){return n.size===void 0?e?e.mergedSize.value:d?.value?.Switch?.size||`medium`:n.size}}),{mergedSizeRef:m,mergedDisabledRef:h}=p,_=r(n.defaultValue),v=e(t(n,`value`),_),y=s(()=>v.value===n.checkedValue),b=r(!1),x=r(!1),S=s(()=>{let{railStyle:e}=n;if(e)return e({focused:x.value,checked:y.value})});function C(e){let{"onUpdate:value":t,onChange:r,onUpdateValue:i}=n,{nTriggerFormInput:a,nTriggerFormChange:o}=p;t&&w(t,e),i&&w(i,e),r&&w(r,e),_.value=e,a(),o()}function T(){let{nTriggerFormFocus:e}=p;e()}function E(){let{nTriggerFormBlur:e}=p;e()}function D(){n.loading||h.value||(v.value===n.checkedValue?C(n.uncheckedValue):C(n.checkedValue))}function O(){x.value=!0,T()}function ee(){x.value=!1,E(),b.value=!1}function te(e){n.loading||h.value||e.key===` `&&(v.value===n.checkedValue?C(n.uncheckedValue):C(n.checkedValue),b.value=!1)}function k(e){n.loading||h.value||e.key===` `&&(e.preventDefault(),b.value=!0)}let A=s(()=>{let{value:e}=m,{self:{opacityDisabled:t,railColor:n,railColorActive:r,buttonBoxShadow:a,buttonColor:o,boxShadowFocus:s,loadingColor:c,textColor:l,iconColor:u,[i(`buttonHeight`,e)]:d,[i(`buttonWidth`,e)]:p,[i(`buttonWidthPressed`,e)]:h,[i(`railHeight`,e)]:g,[i(`railWidth`,e)]:_,[i(`railBorderRadius`,e)]:v,[i(`buttonBorderRadius`,e)]:y},common:{cubicBezierEaseInOut:b}}=f.value,x,S,C;return K?(x=`calc((${g} - ${d}) / 2)`,S=`max(${g}, ${d})`,C=`max(${_}, calc(${_} + ${d} - ${g}))`):(x=M((N(g)-N(d))/2),S=M(Math.max(N(g),N(d))),C=N(g)>N(d)?_:M(N(_)+N(d)-N(g))),{"--n-bezier":b,"--n-button-border-radius":y,"--n-button-box-shadow":a,"--n-button-color":o,"--n-button-width":p,"--n-button-width-pressed":h,"--n-button-height":d,"--n-height":S,"--n-offset":x,"--n-opacity-disabled":t,"--n-rail-border-radius":v,"--n-rail-color":n,"--n-rail-color-active":r,"--n-rail-height":g,"--n-rail-width":_,"--n-width":C,"--n-box-shadow-focus":s,"--n-loading-color":c,"--n-text-color":l,"--n-icon-color":u}}),j=l?u(`switch`,s(()=>m.value[0]),A,n):void 0;return{handleClick:D,handleBlur:ee,handleFocus:O,handleKeyup:te,handleKeydown:k,mergedRailStyle:S,pressed:b,mergedClsPrefix:o,mergedValue:v,checked:y,mergedDisabled:h,cssVars:l?void 0:A,themeClass:j?.themeClass,onRender:j?.onRender}},render(){let{mergedClsPrefix:e,mergedDisabled:t,checked:n,mergedRailStyle:r,onRender:i,$slots:a}=this;i?.();let{checked:o,unchecked:s,icon:c,"checked-icon":l,"unchecked-icon":u}=a,p=!(C(c)&&C(l)&&C(u));return d(`div`,{role:`switch`,"aria-checked":n,class:[`${e}-switch`,this.themeClass,p&&`${e}-switch--icon`,n&&`${e}-switch--active`,t&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},d(`div`,{class:`${e}-switch__rail`,"aria-hidden":`true`,style:r},f(o,t=>f(s,n=>t||n?d(`div`,{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},d(`div`,{class:`${e}-switch__rail-placeholder`},d(`div`,{class:`${e}-switch__button-placeholder`}),t),d(`div`,{class:`${e}-switch__rail-placeholder`},d(`div`,{class:`${e}-switch__button-placeholder`}),n)):null)),d(`div`,{class:`${e}-switch__button`},f(c,t=>f(l,n=>f(u,r=>d(T,null,{default:()=>this.loading?d(m,Object.assign({key:`loading`,clsPrefix:e,strokeWidth:20},this.spinProps)):this.checked&&(n||t)?d(`div`,{class:`${e}-switch__button-icon`,key:n?`checked-icon`:`icon`},n||t):!this.checked&&(r||t)?d(`div`,{class:`${e}-switch__button-icon`,key:r?`unchecked-icon`:`icon`},r||t):null})))),f(o,t=>t&&d(`div`,{key:`checked`,class:`${e}-switch__checked`},t)),f(s,t=>t&&d(`div`,{key:`unchecked`,class:`${e}-switch__unchecked`},t)))))}});export{de as n,fe as t};