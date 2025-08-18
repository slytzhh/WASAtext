<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>
	<div class="container-fluid">
		<RouterView />
	</div>
</template>

<style>
	.container-fluid{
		display: flex;
        justify-content: center;
        align-items: center;
		--bs-gutter-x:0;
		height: 100vh;
	}
</style>