import{En as e,Et as t,Q as n,Vt as r,Yt as i,Z as a,Zt as o,bt as s,c,g as l,hn as u,jt as d,l as f,s as p,wt as m,xt as h}from"./Button-v7qKim5p.js";import{b as g}from"./Tooltip-qMw84P-o.js";import{D as _,p as v}from"./token-dQuugHoF.js";function y(e){let{opacityDisabled:t,heightTiny:n,heightSmall:r,heightMedium:i,heightLarge:a,heightHuge:o,primaryColor:s,fontSize:c}=e;return{fontSize:c,textColor:s,sizeTiny:n,sizeSmall:r,sizeMedium:i,sizeLarge:a,sizeHuge:o,color:s,opacitySpinning:t}}var b={name:`Spin`,common:p,self:y},x=s([s(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),h(`spin-container`,`
 position: relative;
 `,[h(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[v()])]),h(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),h(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[m(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),h(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),h(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[m(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),S={small:20,medium:18,large:16},C=i({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},l.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),f),slots:Object,setup(i){let{mergedClsPrefixRef:o,inlineThemeDisabled:s}=n(i),c=l(`Spin`,`-spin`,x,b,i,o),d=r(()=>{let{size:e}=i,{common:{cubicBezierEaseInOut:n},self:r}=c.value,{opacitySpinning:a,color:o,textColor:s}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof e==`number`?_(e):r[t(`size`,e)],"--n-color":o,"--n-text-color":s}}),f=s?a(`spin`,r(()=>{let{size:e}=i;return typeof e==`number`?String(e):e[0]}),d,i):void 0,p=g(i,[`spinning`,`show`]),m=e(!1);return u(e=>{let t;if(p.value){let{delay:n}=i;if(n){t=window.setTimeout(()=>{m.value=!0},n),e(()=>{clearTimeout(t)});return}}m.value=p.value}),{mergedClsPrefix:o,active:m,mergedStrokeWidth:r(()=>{let{strokeWidth:e}=i;if(e!==void 0)return e;let{size:t}=i;return S[typeof t==`number`?`medium`:t]}),cssVars:s?void 0:d,themeClass:f?.themeClass,onRender:f?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&o(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),s=t.icon?o(`div`,{class:[`${n}-spin-body`,this.themeClass]},o(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):o(`div`,{class:[`${n}-spin-body`,this.themeClass]},o(c,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?o(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},o(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),o(d,{name:`fade-in-transition`},{default:()=>this.active?s:null})):s}});export{C as t};