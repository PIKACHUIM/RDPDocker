import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'deskcli / RDDocker',
  description: '开箱即用的桌面容器环境，支持 Docker/Podman/LXC 多平台',
  lang: 'zh-CN',
  themeConfig: {
    nav: [
      { text: '首页', link: '/' },
      { text: '指南', link: '/guide/introduction' },
      { text: 'API 文档', link: '/api/reference' },
      { text: 'GitHub', link: 'https://github.com/PIKACHUIM/RDDocker' }
    ],
    sidebar: {
      '/guide/': [
        {
          text: '开始',
          items: [
            { text: '项目介绍', link: '/guide/introduction' },
            { text: '快速开始', link: '/guide/quick-start' },
            { text: '环境要求', link: '/guide/requirements' },
            { text: '安装方法', link: '/guide/installation' }
          ]
        },
        {
          text: '使用',
          items: [
            { text: 'Web 管理面板', link: '/guide/management' },
            { text: '镜像使用', link: '/guide/images' },
            { text: 'deskcli CLI', link: '/guide/easydesk' }
          ]
        },
        {
          text: '帮助',
          items: [
            { text: '常见问题', link: '/guide/faq' }
          ]
        }
      ],
      '/api/': [
        { text: 'API 参考', link: '/api/reference' }
      ]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/PIKACHUIM/RDDocker' }
    ],
    footer: {
      message: 'MIT License',
      copyright: 'Copyright © 2024 Pikachu'
    }
  }
})
