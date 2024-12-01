<template>
  <div>
    <h1><FaIcon icon="fa-users" />Users</h1>

    <table class="table table-zebra">
      <thead>
        <tr>
        <td></td>
        <td>Username</td>
        <td>Email</td>
        <td></td>
        </tr>
      </thead>
      <tr v-for="user in users" :key="user.id">
        <td>{{ user.id }}</td>
        <td>{{ user.username }}</td>
        <td>{{ user.email }}</td>
        <td>
          <button class="btn btn-sm">
            <FaIcon icon="fa-trash" @click="askDelete(user)" />
          </button>
          <button class="btn btn-sm">
            <FaIcon icon="fa-pencil" @click="editUser(user)" />
          </button>
          <button class="btn btn-sm">
            <FaIcon icon="fa-key" @click="editPassword(user)" />
          </button>
        </td>
      </tr>
    </table>
    <button class="btn btn-primary btn-sm" @click="createUser()">
      <FaIcon icon="fa-plus" />Add user
    </button>
    <dialog ref="user_delete_dialog" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Delete {{ userToDelete?.username }}?</h3>
        <p class="py-4">
          Do you really want to delete {{ userToDelete?.email }}?
        </p>
        <div class="modal-action">
          <form method="dialog">
            <button class="btn btn-primary" @click="confirmDelete()">
              Yes
            </button>
            <!-- if there is a button in form, it will close the modal -->
            <button class="btn">Cancel</button>
          </form>
        </div>
      </div>
    </dialog>
    <dialog ref="user_create_dialog" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Create user</h3>

        <div class="modal-action">
          <form method="dialog" class="md:container" @submit="confirmCreate">
            <div>
              <label class="label w-full max-w-xs">
                <span class="label-text"
                  ><FaIcon icon="fa-user" />Username</span
                >
              </label>
              <input
                type="text"
                v-model="newUser.username"
                placeholder=""
                class="input input-bordered input-sm w-full max-w-xs"
              />
            </div>
            <div>
              <label class="label w-full max-w-xs">
                <span class="label-text"
                  ><FaIcon icon="fa-envelope" />Email</span
                >
              </label>
              <input
                type="email"
                v-model="newUser.email"
                placeholder=""
                class="input input-bordered input-sm w-full max-w-xs"
              />
            </div>

            <div>
              <label class="label w-full max-w-xs">
                <span class="label-text"><FaIcon icon="fa-key" />Password</span>
              </label>
              <input
                type="password"
                v-model="newUser.password"
                placeholder=""
                class="input input-bordered input-sm w-full max-w-xs"
              />
            </div>

            <div>
              <label class="label w-full max-w-xs">
                <span class="label-text"
                  ><FaIcon icon="fa-key" />Confirm Password</span
                >
              </label>
              <input
                type="password"
                v-model="newUser.passwordConfirm"
                placeholder=""
                class="input input-bordered input-sm w-full max-w-xs"
              />
            </div>
            <button class="btn btn-primary">Create User</button>
          </form>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script lang="ts" setup>
definePageMeta({
  layout: "admin",
});

let users = ref((await usePb().collection("users").getList()).items);

let newUser = {
  username: "",
  email: "",
  password: "",
  passwordConfirm: "",
};

let userToDelete = ref<any>();
let user_delete_dialog = ref<HTMLDialogElement>();
let user_create_dialog = ref<HTMLDialogElement>();
async function askDelete(user: any) {
  userToDelete.value = user;
  user_delete_dialog.value?.showModal();
}
async function confirmDelete() {
  if (userToDelete.value) {
    await usePb().collection("users").delete(userToDelete.value.id);
    userToDelete = ref(null);
    users.value = users.value.filter(
      (u: any) => u.id !== userToDelete.value.id
    );
  }
}

async function createUser() {
  user_create_dialog.value?.showModal();
}

async function confirmCreate(e: Event) {
  e.preventDefault();
  try {
    await usePb().collection("users").create(newUser);
    newUser = {
      username: "",
      email: "",
      password: "",
      passwordConfirm: "",
    };
    users.value.push(newUser as any);
    user_create_dialog.value?.close();
  } catch (e) {}
}

async function editUser(user: any) {}
async function editPassword(user: any) {}
</script>

<style scoped>
span svg {
  margin-right: 0.5rem;
}
</style>
