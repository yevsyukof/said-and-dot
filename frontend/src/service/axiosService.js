import axios from 'axios';

export const axiosInstance = axios.create({
    baseURL: 'http://localhost:3000'    // backend-rest-server URL
})
