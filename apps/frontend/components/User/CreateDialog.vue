<template>
  <dialog ref="user_create_dialog" class="modal">
    <div class="modal-box">
      <h3 class="font-bold text-lg">Create user</h3>

      <div class="modal-action">
        <form
          method="dialog"
          class="md:container"
          @submit="$emit('update:modelValue', newUser)"
        >
          <div>
            <label class="label w-full max-w-xs">
              <span class="label-text"><FaIcon icon="fa-user" />Username</span>
            </label>
            <input
              type="text"
              v-model="newUser?.username"
              placeholder=""
              class="input input-bordered input-sm w-full max-w-xs"
            />
          </div>
          <div>
            <label class="label w-full max-w-xs">
              <span class="label-text"><FaIcon icon="fa-envelope" />Email</span>
            </label>
            <input
              type="email"
              v-model="newUser?.email"
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
              v-model="newUser?.password"
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
              v-model="newUser?.passwordConfirm"
              placeholder=""
              class="input input-bordered input-sm w-full max-w-xs"
            />
          </div>
          <button class="btn btn-primary">Create User</button>
        </form>
      </div>
    </div>
  </dialog>
</template>

<script lang="ts" setup>
const props = defineProps({
  modelValue: Object,
});

const emits = defineEmits(["update:modelValue"]);
const newUser = props.modelValue;
</script>
