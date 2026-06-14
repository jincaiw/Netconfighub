import{Ct as e,Dt as t,Et as n,Ot as r,Q as i,Tt as a,Vt as o,Yt as s,Z as c,Zt as l,_t as u,bt as d,g as f,s as p,wt as m,xt as h}from"./Button-v7qKim5p.js";import{b as g,l as _}from"./Tooltip-qMw84P-o.js";import{C as v}from"./token-dQuugHoF.js";import{g as y}from"./index-CESIsNoN.js";function b(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var x={thPaddingBorderedSmall:`8px 12px`,thPaddingBorderedMedium:`12px 16px`,thPaddingBorderedLarge:`16px 24px`,thPaddingSmall:`0`,thPaddingMedium:`0`,thPaddingLarge:`0`,tdPaddingBorderedSmall:`8px 12px`,tdPaddingBorderedMedium:`12px 16px`,tdPaddingBorderedLarge:`16px 24px`,tdPaddingSmall:`0 0 8px 0`,tdPaddingMedium:`0 0 12px 0`,tdPaddingLarge:`0 0 16px 0`};function S(e){let{tableHeaderColor:t,textColor2:n,textColor1:r,cardColor:i,modalColor:a,popoverColor:o,dividerColor:s,borderRadius:c,fontWeightStrong:l,lineHeight:d,fontSizeSmall:f,fontSizeMedium:p,fontSizeLarge:m}=e;return Object.assign(Object.assign({},x),{lineHeight:d,fontSizeSmall:f,fontSizeMedium:p,fontSizeLarge:m,titleTextColor:r,thColor:u(i,t),thColorModal:u(a,t),thColorPopover:u(o,t),thTextColor:r,thFontWeight:l,tdTextColor:n,tdColor:i,tdColorModal:a,tdColorPopover:o,borderColor:u(i,s),borderColorModal:u(a,s),borderColorPopover:u(o,s),borderRadius:c})}var C={name:`Descriptions`,common:p,self:S},w=d([h(`descriptions`,{fontSize:`var(--n-font-size)`},[h(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),h(`descriptions-table-wrapper`,[h(`descriptions-table`,[h(`descriptions-table-row`,[h(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),h(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),a(`bordered`,[h(`descriptions-table-wrapper`,[h(`descriptions-table`,[h(`descriptions-table-row`,[d(`&:last-child`,[h(`descriptions-table-content`,{paddingBottom:0})])])])])]),m(`left-label-placement`,[h(`descriptions-table-content`,[d(`> *`,{verticalAlign:`top`})])]),m(`left-label-align`,[d(`th`,{textAlign:`left`})]),m(`center-label-align`,[d(`th`,{textAlign:`center`})]),m(`right-label-align`,[d(`th`,{textAlign:`right`})]),m(`bordered`,[h(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[h(`descriptions-table`,[h(`descriptions-table-row`,[d(`&:not(:last-child)`,[h(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),h(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),h(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[d(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),h(`descriptions-table-content`,[d(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),h(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),h(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[h(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[h(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[h(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),h(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[e(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),e(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),t(h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),r(h(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),T=`DESCRIPTION_ITEM_FLAG`;function E(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var D=s({name:`Descriptions`,props:Object.assign(Object.assign({},f.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:r,mergedComponentPropsRef:a}=i(e),s=o(()=>e.size||a?.value?.Descriptions?.size||`medium`),l=f(`Descriptions`,`-descriptions`,w,C,e,t),u=o(()=>{let{bordered:t}=e,r=s.value,{common:{cubicBezierEaseInOut:i},self:{titleTextColor:a,thColor:o,thColorModal:c,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[n(`fontSize`,r)]:S,[n(t?`thPaddingBordered`:`thPadding`,r)]:C,[n(t?`tdPaddingBordered`:`tdPadding`,r)]:w}}=l.value;return{"--n-title-text-color":a,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":i,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":o,"--n-th-color-modal":c,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),d=r?c(`descriptions`,o(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=s.value[0],t}),u,e):void 0;return{mergedClsPrefix:t,cssVars:r?void 0:u,themeClass:d?.themeClass,onRender:d?.onRender,compitableColumn:g(e,[`columns`,`column`]),inlineThemeDisabled:r,mergedSize:s}},render(){let e=this.$slots.default,t=e?y(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:i,labelPlacement:a,labelAlign:o,mergedSize:s,bordered:c,title:u,cssVars:d,mergedClsPrefix:f,separator:p,onRender:m}=this;m?.();let h=t.filter(e=>E(e)),g=h.reduce((e,t,o)=>{let s=t.props||{},u=h.length-1===o,d=[`label`in s?s.label:b(t,`label`)],m=[b(t)],g=s.span||1,_=e.span;e.span+=g;let v=s.labelStyle||s[`label-style`]||this.labelStyle,y=s.contentStyle||s[`content-style`]||this.contentStyle;if(a===`left`)c?e.row.push(l(`th`,{class:[`${f}-descriptions-table-header`,r],colspan:1,style:v},d),l(`td`,{class:[`${f}-descriptions-table-content`,n],colspan:u?(i-_)*2+1:g*2-1,style:y},m)):e.row.push(l(`td`,{class:`${f}-descriptions-table-content`,colspan:u?(i-_)*2:g*2},l(`span`,{class:[`${f}-descriptions-table-content__label`,r],style:v},[...d,p&&l(`span`,{class:`${f}-descriptions-separator`},p)]),l(`span`,{class:[`${f}-descriptions-table-content__content`,n],style:y},m)));else{let t=u?(i-_)*2:g*2;e.row.push(l(`th`,{class:[`${f}-descriptions-table-header`,r],colspan:t,style:v},d)),e.secondRow.push(l(`td`,{class:[`${f}-descriptions-table-content`,n],colspan:t,style:y},m))}return(e.span>=i||u)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),a!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>l(`tr`,{class:`${f}-descriptions-table-row`},e));return l(`div`,{style:d,class:[`${f}-descriptions`,this.themeClass,`${f}-descriptions--${a}-label-placement`,`${f}-descriptions--${o}-label-align`,`${f}-descriptions--${s}-size`,c&&`${f}-descriptions--bordered`]},u||this.$slots.header?l(`div`,{class:`${f}-descriptions-header`},u||_(this,`header`)):null,l(`div`,{class:`${f}-descriptions-table-wrapper`},l(`table`,{class:`${f}-descriptions-table`},l(`tbody`,null,a===`top`&&l(`tr`,{class:`${f}-descriptions-table-row`,style:{visibility:`collapse`}},v(i*2,l(`td`,null))),g))))}}),O={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},k=s({name:`DescriptionsItem`,[T]:!0,props:O,slots:Object,render(){return null}});export{D as n,k as t};