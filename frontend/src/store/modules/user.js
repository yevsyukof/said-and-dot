export default {
    state: () => ({
        id: '', // поменял с _id
        firstName: '',
        lastName: '',
        username: '',
        email: '',
        avatar: "https://w7.pngwing.com/pngs/419/473/png-transparent-computer-icons-user-profile-login-user-heroes-sphere-black-thumbnail.png",
        isAdmin: '',
        isMod: '',
        followers: [],
        follows: []
    }),
    getters: {
        getUser() {
            return {
                id: state.id,  // поменял с _id
                firstName: state.firstName,
                lastName: state.lastName,
                username: state.username,
                email: state.email,
                avatar: state.avatar,
                isAdmin: state.isAdmin,
                isMod: state.isMod,
                followers: [],
                follows: []
            }
        }
    },
    mutations: {
        SET_USER(curState, payload) {  // поменял state на curState
            curState.id = payload.id;
            curState.firstname = payload.firstName;
            curState.lastname = payload.lastName;
            curState.username = payload.username;
            curState.email = payload.email;
            // curState.avatar = payload.avatar;
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