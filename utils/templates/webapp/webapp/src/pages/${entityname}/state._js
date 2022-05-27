import {reactive} from "vue";

// export const url = "http://localhost:8081/{{LowerCase .EntityName}}"

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