<script>
import NavBar from './components/NavBar.vue';
import ErrorToast from './components/ErrorToast.vue';

export default {
	components: {
		"nav-bar": NavBar,
		"error-toast": ErrorToast
	}
};
</script>

<template>
  <div>
    <nav-bar/>
    <error-toast/>
    <div class="container-fluid">
      <br>
      <router-view></router-view>
    </div>
  </div>
</template>
