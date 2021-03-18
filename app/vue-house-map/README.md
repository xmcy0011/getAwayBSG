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

### how to use elementui in vue3

vue中如何使用element?
```bash
vue add element # vue-cli use element-ui
```

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