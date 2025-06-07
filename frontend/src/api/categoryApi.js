import axios from 'axios';

const baseUrl = "http://localhost:8080";

export const listCategories = async (page = 1, limit = 10) => {
  try {
    const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
    if (!token) {
      throw new Error("No authentication token found. Please log in.");
    }
    const response = await axios.get(`${baseUrl}/api/categories`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      params: {
        page,
        limit,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching categories:', error);
    throw error;
  }
}; 