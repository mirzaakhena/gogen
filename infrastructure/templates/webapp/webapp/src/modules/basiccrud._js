import {reactive, toRefs} from "vue";
import to from "await-to-js";
import axios from "axios";

export default function BasicCrud(url) {

    const loadData = async (params) => {

        const requestConfig = { params: { ...params } }

        const [err, res] = await to(axios.get(url, requestConfig).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)

    }

    const addNewData = async (payload) => {

        const [err, res] = await to(axios.post(url, payload).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)
    }

    const deleteData = async (id) => {

        const [err, res] = await to(axios.delete(`${url}/${id}`).catch((err) => Promise.reject(err)))

        if (err) {
            return Promise.reject(err.response.data)
        }

        return Promise.resolve(res.data)
    }

    const updateData = async (id, payload) => {

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