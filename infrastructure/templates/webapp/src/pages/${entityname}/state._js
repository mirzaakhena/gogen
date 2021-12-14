import {reactive} from "vue";

export const state = reactive({
    items: [],
    item: {},
    filter: {
        page: 1,
        size: 4,
    },
    totalItems: 0,
})

export const getNumberOfPage = () => {
    return Math.ceil(state.totalItems / state.filter.size)
}