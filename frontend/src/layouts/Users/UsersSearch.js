import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import IconButton from "@mui/material/IconButton";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import Alert from "@mui/material/Alert";
import MDBox from "components/MDBox";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import DataTable from "examples/Tables/DataTable";
import MDPagination from "components/MDPagination";
import ConfirmDialog from "components/ConfirmDialog";
import userApi from '../../api/userApi';

const UsersSearch = () => {
  const [searchResults, setSearchResults] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [alert, setAlert] = useState(null);
  const navigate = useNavigate();
  const [openConfirm, setOpenConfirm] = useState(false);
  const [userIdToDelete, setUserIdToDelete] = useState(null);

  useEffect(() => {
    const fetchUsersData = async () => {
      try {
        const data = await userApi.fetchUsers(currentPage, rowsPerPage);
        setSearchResults(data.data);
        setTotalPages(data.pagination.total_pages);
      } catch (error) {
        setAlert({
          message: "Erro ao buscar usuários: " + error.message,
          color: "error",
        });
      }
    };
    fetchUsersData();
  }, [currentPage, rowsPerPage]);

  const handleEdit = (user) => {
    navigate("/UpdateUser", { state: { userData: user, userId: user.id } });
  };

  const askDelete = (userId) => {
    setUserIdToDelete(userId);
    setOpenConfirm(true);
  };

  const handleDelete = async (userId) => {
    try {
      await axios.delete(`http://localhost:3000/api/users/${userId}`);
      setSearchResults(searchResults.filter((user) => user.id !== userId));
      setAlert({
        message: "Usuário deletado com sucesso.",
        color: "success",
      });
    } catch (error) {
      setAlert({
        message: "Erro ao apagar usuário: " + error.message,
        color: "error",
      });
    }
  };

  const confirmDelete = async () => {
    setOpenConfirm(false);
    if (userIdToDelete) {
      await handleDelete(userIdToDelete);
      setUserIdToDelete(null);
    }
  };

  const handlePageChange = (event, value) => {
    setCurrentPage(value);
  };

  const paginatedResults = Array.isArray(searchResults)
    ? searchResults.slice((currentPage - 1) * rowsPerPage, currentPage * rowsPerPage)
    : [];

  return (
    <DashboardLayout>
      <DashboardNavbar showTitle={false} />
      <MDBox pt={6} pb={3}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <MDBox
                mx={2}
                mt={-3}
                py={3}
                px={2}
                variant="gradient"
                bgColor="info"
                borderRadius="lg"
                coloredShadow="info"
              >
                <h6 style={{ color: "white", margin: 0 }}>Usuários</h6>
              </MDBox>
              <MDBox p={3}>
                {alert && (
                  <Alert
                    severity={alert.color === "error" ? "error" : "success"}
                    onClose={() => setAlert(null)}
                  >
                    {alert.message}
                  </Alert>
                )}
                <DataTable
                  table={{
                    columns: [
                      { Header: "ID", accessor: "id", width: "5%" },
                      { Header: "Username", accessor: "username", width: "20%" },
                      { Header: "Email", accessor: "email", width: "20%" },
                      { Header: "Active", accessor: "active", width: "10%" },
                      { Header: "Created At", accessor: "created_at", width: "20%" },
                      { Header: "Updated At", accessor: "updated_at", width: "20%" },
                      // {
                      //   Header: "Ações",
                      //   accessor: "actions",
                      //   width: "5%",
                      //   Cell: ({ row }) => (
                      //     <div style={{ display: "flex", gap: "8px" }}>
                      //       <IconButton color="primary" onClick={() => handleEdit(row.original)}>
                      //         <EditIcon />
                      //       </IconButton>
                      //       <IconButton color="secondary" onClick={() => askDelete(row.original.id)}>
                      //         <DeleteIcon />
                      //       </IconButton>
                      //     </div>
                      //   ),
                      // },
                    ],
                    rows: paginatedResults,
                  }}
                  page={currentPage}
                />
                <MDBox mt={3} display="flex" justifyContent="center">
                  <MDPagination
                    count={totalPages}
                    page={currentPage}
                    onChange={handlePageChange}
                  />
                </MDBox>
              </MDBox>
            </Card>
          </Grid>
        </Grid>
      </MDBox>
      <ConfirmDialog
        open={openConfirm}
        title="Confirmar exclusão"
        message="Tem certeza que deseja deletar este usuário?"
        onClose={() => setOpenConfirm(false)}
        onConfirm={confirmDelete}
      />
      <Footer />
    </DashboardLayout>
  );
};

export default UsersSearch; 