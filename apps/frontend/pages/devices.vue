<template>
  <div>
    <h1>Devices</h1>
    <button class="btn btn-primary" @click="linkDevice">
      <FaIcon icon="computer" /> Link Device
    </button>
    <table class="table">
      <tbody>
        <tr v-for="device in devices" :key="device.id">
          <td>
            <FaIcon
              icon="check-circle"
              class="mr-2"
              :class="device.verified ? 'text-success' : ''"
            />{{ device.token }}
          </td>
          <td>{{ device.name }}</td>
          <td>{{ device.created }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts" setup>
const devices = ref();
async function generateID() {
  const id = Math.random().toString(36).substring(2).toUpperCase();
  const existing = await usePb()
    .collection("devices")
    .getFullList({ filter: `(token='${id}')` });
  if (existing.length > 0) {
    return generateID();
  } else {
    return id;
  }
}
async function linkDevice() {
  await usePb()
    .collection("devices")
    .create({
      user: usePb().authStore.model?.id,
      token: await generateID(),
      verified: false,
      name: "My Device",
    });
  getDevices();
}
async function getDevices() {
  devices.value = await usePb()
    .collection("devices")
    .getFullList({ filter: `(user='${usePb().authStore.model?.id}')` });
}
getDevices();
</script>
