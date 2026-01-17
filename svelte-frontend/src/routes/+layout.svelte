 <script lang="ts">
 	import '../app.css';

 	let { children } = $props();
 </script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
	<nav class="bg-white/80 backdrop-blur-sm shadow-sm border-b border-white/20 sticky top-0 z-50">
		<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
 			<div class="flex justify-between items-center h-16">
 				<a href="/" class="text-xl sm:text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">
 					My Journal
 				</a>
 				<a href="/new" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">New Entry</a>
 			</div>
		</div>
	</nav>
	<main class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{@render children()}
	</main>
</div>