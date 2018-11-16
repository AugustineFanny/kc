import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/Login'
import Navbar from '@/components/Navbar'
import Users from '@/components/Users'
import Wallets from '@/components/Wallets'
import Addresses from '@/components/Addresses'
import RealNames from '@/components/RealNames'
import Kycs from '@/components/Kycs'
import TransferOuts from '@/components/TransferOuts'
import Transfers from '@/components/Transfers'
import Ads from '@/components/Ads'
import Trades from '@/components/Trades'
import Chats from '@/components/Chats'
import TradeLogs from '@/components/TradeLogs'
import Datums from '@/components/Datums'
import Appeals from '@/components/Appeals'
import Admins from '@/components/Admins'
import AddressPool from '@/components/AddressPool'
import Logs from '@/components/Logs'
import Coins from '@/components/Coins'
import Statistics from '@/components/Statistics'
import Invites from '@/components/Invites'
import ChildAmounts from '@/components/ChildAmounts'
import Groups from '@/components/Groups'
import FundChanges from '@/components/FundChanges'
import Announcements from '@/components/Announcements'
import BatchSms from '@/components/BatchSms'
import BatchEmail from '@/components/BatchEmail'
import Online from '@/components/Online'
import Banners from '@/components/Banners'
import LedgerBTC from '@/components/LedgerBTC'
import LedgerETH from '@/components/LedgerETH'
import ProfitDate from '@/components/ProfitDate'
import ProfitMonth from '@/components/ProfitMonth'
import Predistribution from '@/components/Predistribution'
import Subscription from '@/components/Subscription'

import Simulation from '@/components/Simulation'
import SUser from '@/components/SUser'
import SActivity from '@/components/SActivity'

import BtcdoScript from '@/components/BtcdoScript'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '',
      redirect: '/user'
    }, 
    {
      path: '/login',
      component: Login
    },
    { path: '/simulation', component: Simulation},
    { path: '/simulation/user', component: SUser},
    { path: '/simulation/activity', component: SActivity},
    { path: '/btcdoscript', component: BtcdoScript},
    {
      path: '/user',
      component: Navbar,
      children: [
        { path: '', component: Users},
        { path: 'real-name/:uid', component: RealNames},
        { path: 'kyc/:uid', component: Kycs},
        { path: 'wallets', component: Wallets},
        { path: 'addresses', component: Addresses},
        { path: 'real-names', component: RealNames},
        { path: 'kycs', component: Kycs},
        { path: 'transfer-outs', component: TransferOuts},
        { path: 'transfers', component: Transfers},
        { path: 'fund-changes', component: FundChanges}
      ]
    },
    {
      path: '/trades',
      component: Navbar,
      children: [
        { path: '', component: Trades},
      ]
    },
    {
      path: '/trade/:code/',
      component: Navbar,
      children: [
        { path: 'chats', component: Chats},
        { path: 'logs', component: TradeLogs},
        { path: 'datums', component: Datums},
      ]
    },
    {
      path: '/stat',
      component: Navbar,
      children: [
        { path: 'statistics', component: Statistics},
        { path: 'invites', component: Invites},
        { path: 'child/amounts', component: ChildAmounts}
      ]
    },
    {
      path: '/manage',
      component: Navbar,
      children: [
        { path: 'subscription', component: Subscription},
        { path: 'simulation', redirect: '/simulation'},
        { path: 'predistribution', component: Predistribution},
        { path: 'profit-date', component: ProfitDate},
        { path: 'profit-month', component: ProfitMonth}
      ]
    },
    {
      path: '/conf',
      component: Navbar,
      children: [
        { path: 'banners', component: Banners},
        { path: 'announcements', component: Announcements}
      ]
    },
    {
      path: '/admin',
      component: Navbar,
      children: [
        { path: 'admins', component: Admins},
        { path: 'coins', component: Coins},
        { path: 'groups', component: Groups},
        { path: 'address-pool', component: AddressPool},
        { path: 'logs', component: Logs},
        { path: 'batch-sms', component: BatchSms},
        { path: 'batch-email', component: BatchEmail},
        { path: 'online', component: Online},
        { path: 'ledgerBTC', component: LedgerBTC},
        { path: 'ledgerETH', component: LedgerETH},
      ]
    }
  ]
})
