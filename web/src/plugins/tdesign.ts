// On-demand TDesign registration. Only components actually used by the app are
// registered and only their styles are bundled (plus the small global token
// stylesheet, which defines the --td-* CSS variables every component needs).
import type { App, Component } from 'vue'
// TDesign 1.20.11 has no Result export/style; use a local equivalent for <t-result>.
import TResult from '@/components/TResult.vue'
import {
  Alert,
  Aside,
  Badge,
  Breadcrumb,
  BreadcrumbItem,
  Button,
  Card,
  Col,
  ConfigProvider,
  Content,
  DatePicker,
  Descriptions,
  DescriptionsItem,
  Dialog,
  Divider,
  Drawer,
  Dropdown,
  DropdownItem,
  Empty,
  Form,
  FormItem,
  Header,
  Input,
  InputNumber,
  Layout,
  Loading,
  Menu,
  MenuItem,
  Radio,
  RadioGroup,
  Row,
  Select,
  Space,
  Switch,
  TabPanel,
  Table,
  Tabs,
  Tag,
  Textarea,
  Upload,
} from 'tdesign-vue-next'

import 'tdesign-vue-next/es/style/index.css'
import 'tdesign-vue-next/es/alert/style/index.css'
import 'tdesign-vue-next/es/badge/style/index.css'
import 'tdesign-vue-next/es/breadcrumb/style/index.css'
import 'tdesign-vue-next/es/button/style/index.css'
import 'tdesign-vue-next/es/card/style/index.css'
import 'tdesign-vue-next/es/date-picker/style/index.css'
import 'tdesign-vue-next/es/descriptions/style/index.css'
import 'tdesign-vue-next/es/dialog/style/index.css'
import 'tdesign-vue-next/es/divider/style/index.css'
import 'tdesign-vue-next/es/drawer/style/index.css'
import 'tdesign-vue-next/es/dropdown/style/index.css'
import 'tdesign-vue-next/es/empty/style/index.css'
import 'tdesign-vue-next/es/form/style/index.css'
import 'tdesign-vue-next/es/grid/style/index.css'
import 'tdesign-vue-next/es/input/style/index.css'
import 'tdesign-vue-next/es/input-number/style/index.css'
import 'tdesign-vue-next/es/layout/style/index.css'
import 'tdesign-vue-next/es/loading/style/index.css'
import 'tdesign-vue-next/es/menu/style/index.css'
import 'tdesign-vue-next/es/message/style/index.css'
import 'tdesign-vue-next/es/notification/style/index.css'
import 'tdesign-vue-next/es/radio/style/index.css'
import 'tdesign-vue-next/es/select/style/index.css'
import 'tdesign-vue-next/es/space/style/index.css'
import 'tdesign-vue-next/es/switch/style/index.css'
import 'tdesign-vue-next/es/table/style/index.css'
import 'tdesign-vue-next/es/tabs/style/index.css'
import 'tdesign-vue-next/es/tag/style/index.css'
import 'tdesign-vue-next/es/textarea/style/index.css'
import 'tdesign-vue-next/es/upload/style/index.css'

const components: Record<string, Component> = {
  TAlert: Alert,
  TAside: Aside,
  TBadge: Badge,
  TBreadcrumb: Breadcrumb,
  TBreadcrumbItem: BreadcrumbItem,
  TButton: Button,
  TCard: Card,
  TCol: Col,
  TConfigProvider: ConfigProvider,
  TContent: Content,
  TDatePicker: DatePicker,
  TDescriptions: Descriptions,
  TDescriptionsItem: DescriptionsItem,
  TDialog: Dialog,
  TDivider: Divider,
  TDrawer: Drawer,
  TDropdown: Dropdown,
  TDropdownItem: DropdownItem,
  TEmpty: Empty,
  TForm: Form,
  TFormItem: FormItem,
  THeader: Header,
  TInput: Input,
  TInputNumber: InputNumber,
  TLayout: Layout,
  TLoading: Loading,
  TMenu: Menu,
  TMenuItem: MenuItem,
  TRadio: Radio,
  TRadioGroup: RadioGroup,
  TRow: Row,
  TResult,
  TSelect: Select,
  TSpace: Space,
  TSwitch: Switch,
  TTabPanel: TabPanel,
  TTable: Table,
  TTabs: Tabs,
  TTag: Tag,
  TTextarea: Textarea,
  TUpload: Upload,
}

export function setupTDesign(app: App): void {
  for (const [name, component] of Object.entries(components)) {
    app.component(name, component)
  }
}
