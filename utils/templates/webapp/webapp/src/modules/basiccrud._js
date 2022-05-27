import to from "await-to-js";
import axios from "axios";
import {BASE_URL} from "./url";

export default function BasicCrud() {

    const loadData = async (params) => {

        const url = `${BASE_URL}/{{LowerCase .EntityName}}`

        const requestConfig = { params: { ...params } }

        const [err, res] = await to(axios.get(url, requestConfig).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)

    }

    const addNewData = async (payload) => {

        const url = `${BASE_URL}/{{LowerCase .EntityName}}`

        const [err, res] = await to(axios.post(url, payload).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)
    }

    const deleteData = async (id) => {

        const url = `${BASE_URL}/{{LowerCase .EntityName}}`

        const [err, res] = await to(axios.delete(`${url}/${id}`).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)
    }

    const updateData = async (id, payload) => {

        const url = `${BASE_URL}/{{LowerCase .EntityName}}`

        const [err, res] = await to(axios.put(`${url}/${id}`, payload).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)
    }

    return {
        loadData,
        addNewData,
        deleteData,
        updateData,
    }

}