# house-map

## introduce

use **vue3** , **element-ui-3.x** and **electron11(chrome89)** build cross-platform desktop client.

## quick start

## install

1. npm
```

```

2. vue
```
```

3. vue-cli
```
```

3. element-ui-3.x
```bash
$ cnpm install element-plus --save # install element-ui-3.x
```

4. nodejs
```
```

5. electron
```
$ 
$ yarn add vue-cli-plugin-electron-builder -D --tilde
```

## build

```bash
$ yarn install    # install depends
$ yarn serve      # Compiles and hot-reloads for development
$ yarn build      # Compiles and minifies for production
$ yarn lint       # Lints and fixes files
```

## run

run in browser:
```
$ yarn serve
```

run in desktop:
```
$ yarn electron:serve
```

## Note

### how to use electron in vue3

如何在vue3中使用electron？
参考：[使用 vue-cli3.x 快速构建 electron 项目](https://www.icode9.com/content-4-422756.html)

```bash
$ vue create my-project     # 使用 Vue CLI3 创建项目vue3项目("my-project" 改为自己想要取是项目名称)
$ cd my-project
$ vue add electron-builder  # vue3中使用electron
$ npm run electron:serve    # 开发模式运行
$ npm run electron:build    # 部署，编译
```

### how to use gaode map in vue3

1. use [JSAPI Loader](https://lbs.amap.com/api/jsapi-v2/guide/abc/load)
2. [Vue-Cli 3.0 中配置高德地图api](https://blog.csdn.net/weixin_43953710/article/details/101377497)

### how to use elementui in vue3

vue中如何使用element?
```bash
vue add element # vue-cli use element-ui
```

### how to use axios in vue3

see:
- [vue3.0引入axios](https://blog.csdn.net/u011724770/article/details/114403738)
- [vue3.0使用axios-二次封装](https://blog.csdn.net/np918/article/details/95018440)

1. mian.js
```js
import axios from 'axios'
import VueAxios from 'vue-axios'
 
createApp(App).use(VueAxios, axios)
```

2. 使用
```js
vue.axios.get(api).then((response) => {
    console.log(response.data)
})
this.axios.get(api).then((response) => { 
    console.log (response.data)
})
this.$http.get(api).then((response) =>{
    console.log( response.data)
})
```

### async and await

see:[Javascript中async与await的用法](https://blog.csdn.net/weixin_42042017/article/details/109472908)

### install depends by manual

```bash
$ cnpm install element-plus --save         # element-ui
$ cnpm install babel-plugin-import --save  # 
```

## 常见问题

### 安装electron失败 postinstall: `node install.js`

解决方法：将electron下载地址指向taobao镜像
```bash
npm config set electron_mirror "https://npm.taobao.org/mirrors/electron/"
```

### Error: Rule can only have one resource source (provided resource and test + include + exclude)

package.json中webpack版本冲突问题。删除webpack，重新装以前的版本。
```bash
npm uninstall webpack
npm install webpack@^4.0.0 --save-dev
```