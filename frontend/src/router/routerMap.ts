/**
 * 基础路由
 * @type { *[] }
 */

const constantRouterMap = [
  {
    path: '/',
    component: () => import('@/layouts/AppSider.vue'),
    redirect: { name: 'connect' },
    children: [
      {
        path: '/connect',
        name: 'connect',
        component: () => import('@/views/connect/index.vue'),
        params: { id: 'connect' },
        meta: {
          title: '连接',
          locale: 'menu.connect',
          icon: 'icon-connect',
          selecteicon: 'icon-connect-selecte',
        }
      },
      {
        path: '/stress',
        name: 'stress',
        component: () => import('@/views/stress/index.vue'),
        params: { id: 'connect' },
        meta: {
          title: '压测',
          locale: 'menu.stress',
          icon: 'icon-stress',
          selecteicon: 'icon-stress-selecte',
        }
      },
      {
        path: '/gmqt',
        name: 'gmqt',
        component: () => import('@/views/gmqt/index.vue'),
        params: { id: 'gmqt' },
        meta: {
          title: 'GMQT',
          locale: 'menu.gmqt',
          icon: 'icon-MQTT2',
          selecteicon: 'icon-MQTT2',
        }
      },
      {
        path: '/help',
        name: 'help',
        component: () => import('@/views/help/index.vue'),
        params: { id: 'help' },
        meta: {
          title: '关于',
          locale: 'menu.help',
          icon: 'icon-guanyu2',
          selecteicon: 'icon-guanyu2',
        }
      },
      {
        path: '/setting',
        name: 'setting',
        component: () => import('@/views/setting/index.vue'),
        meta: {
          hideInMenu:true,
        }
      },      
    ]
  },
]

export default constantRouterMap