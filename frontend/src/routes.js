// Material Dashboard 2 React layouts
import Dashboard from "layouts/dashboard";
import Tables from "layouts/tables";
import Billing from "layouts/billing";
import RTL from "layouts/rtl";
import Notifications from "layouts/notifications";
import Profile from "layouts/profile";
import SignIn from "layouts/authentication/sign-in";
import SignUp from "layouts/authentication/sign-up";
import ProtectedRoute from "components/ProtectedRoute";
import UsersSearch from "layouts/Users/UsersSearch";
import ImportarSearch from "layouts/importar/importarSearch";
import ImportarAdd from "layouts/importar/importarAdd";

// @mui icons
import Icon from "@mui/material/Icon";

const routes = [
  {
    type: "collapse",
    name: "Dashboard",
    key: "dashboard",
    icon: <Icon fontSize="small">dashboard</Icon>,
    route: "/dashboard",
    component: (
      <ProtectedRoute>
        <Dashboard />
      </ProtectedRoute>
    ),
  },
  // {
  //   type: "collapse",
  //   name: "Tables",
  //   key: "tables",
  //   icon: <Icon fontSize="small">table_view</Icon>,
  //   route: "/tables",
  //   component: (
  //     <ProtectedRoute>
  //       <Tables />
  //     </ProtectedRoute>
  //   ),
  // },
  // {
  //   type: "collapse",
  //   name: "Billing",
  //   key: "billing",
  //   icon: <Icon fontSize="small">receipt_long</Icon>,
  //   route: "/billing",
  //   component: (
  //     <ProtectedRoute>
  //       <Billing />
  //     </ProtectedRoute>
  //   ),
  // },
  // {
  //   type: "collapse",
  //   name: "RTL",
  //   key: "rtl",
  //   icon: <Icon fontSize="small">format_textdirection_r_to_l</Icon>,
  //   route: "/rtl",
  //   component: (
  //     <ProtectedRoute>
  //       <RTL />
  //     </ProtectedRoute>
  //   ),
  // },
  // {
  //   type: "collapse",
  //   name: "Notifications",
  //   key: "notifications",
  //   icon: <Icon fontSize="small">notifications</Icon>,
  //   route: "/notifications",
  //   component: (
  //     <ProtectedRoute>
  //       <Notifications />
  //     </ProtectedRoute>
  //   ),
  // },
  // {
  //   type: "collapse",
  //   name: "Profile",
  //   key: "profile",
  //   icon: <Icon fontSize="small">person</Icon>,
  //   route: "/profile",
  //   component: (
  //     <ProtectedRoute>
  //       <Profile />
  //     </ProtectedRoute>
  //   ),
  // },

  {
    type: "collapse",
     name: "Usuarios",
     key: "usuarios",
     icon: <Icon fontSize="small">group</Icon>,
     route: "/usuarios",
     component: <UsersSearch />,
   },
     {
     type: "collapse",
     name: "Adicionar Importação",
     key: "adicionar-importacao",
     icon: <Icon fontSize="small">add</Icon>,
     route: "/importar/add",
     component: <ImportarAdd />,
   },
   {
     type: "collapse",
     name: "Buscar Importações",
     key: "importacoes",
     icon: <Icon fontSize="small">cloud_upload</Icon>,
     route: "/importar",
     component: <ImportarSearch />,
   },
 
   {
     type: "collapse",
     name: "Sign In",
    key: "sign-in",
     icon: <Icon fontSize="small">login</Icon>,
     route: "/authentication/sign-in",
     component: <SignIn />,
   },
   {
     type: "collapse",
     name: "Sign Up",
     key: "sign-up",
     icon: <Icon fontSize="small">assignment</Icon>,
     route: "/authentication/sign-up",
     component: <SignUp />,
   },
 
];

export default routes;
