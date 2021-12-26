import to from "await-to-js";
import axios from "axios";

export const orderSubmitted = async (payload) => {

    const [err, res] = await to(axios.post("http://localhost:8081/order", payload).catch((err) => Promise.reject(err)))

    if (err) {
        return Promise.reject(err.response.data)
    }

    return Promise.resolve(res.data)
}