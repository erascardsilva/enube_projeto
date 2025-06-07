const authService = {
  logout: () => {
    localStorage.removeItem("user");
  },

  getCurrentUser: () => {
    return JSON.parse(localStorage.getItem("user"));
  },

  getToken: () => {
    const user = JSON.parse(localStorage.getItem("user"));
    return user?.token;
  },

  isAuthenticated: () => {
    const user = JSON.parse(localStorage.getItem("user"));
    return !!user?.token;
  },
};

export default authService; 