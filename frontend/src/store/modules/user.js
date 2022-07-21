export default {
    state: () => ({
        _id: '',
        firstname: '',
        lastname: '',
        username: '',
        email: '',
        avatar: '',
        isAdmin: '',
        isMod: '',
    }),
    getters: {
        getUser() {
            return user = {
                _id: state.id,
                firstname: state.firstname,
                lastname: state.lastname,
                username: state.username,
                email: state.email,
                avatar: state.avatar,
                isAdmin: state.isAdmin,
                isMod: state.isMod,
            }
        }
    },
    mutations: {
        SET_USER(curState, payload) {  // поменял state на curState
            curState._id = payload._id;
            curState.firstname = payload.firstname;
            curState.lastname = payload.lastname;
            curState.username = payload.username;
            curState.email = payload.email;
            curState.avatar = payload.avatar;
            curState.isAdmin = payload.isAdmin;
            curState.isMod = payload.isMod;
        }
    },
    actions: {
        saveUser({commit}, data) {
            commit('SET_USER', data)
        }
    }
}