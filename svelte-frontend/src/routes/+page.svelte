 <script lang="ts">
 	import { onMount } from 'svelte';
 	import { journalApi } from '$lib/api/journal';
 	import type { JournalEntry } from '$lib/types/journal';

	let entries: JournalEntry[] = [];
	let loading = true;
	let error: string | null = null;
	let searchTerm = '';

	onMount(async () => {
		await loadEntries();
	});

	async function loadEntries() {
		try {
			error = null;
			entries = await journalApi.getEntries();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load entries';
		} finally {
			loading = false;
		}
	}

  async function deleteEntry(id: number) {
    if (!confirm('Are you sure you want to delete this entry?')) {
      return;
    }

    try {
      await journalApi.deleteEntry(id);
      entries = entries.filter(entry => entry.id !== id);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete entry';
    }
  }

	$: filteredEntries = entries.filter(entry =>
		entry.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
		entry.content.toLowerCase().includes(searchTerm.toLowerCase())
	);

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<div class="space-y-8">
 	<div class="flex flex-col sm:flex-row sm:justify-between sm:items-start space-y-6 sm:space-y-0">
 		<div>
 			<h1 class="text-3xl sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600 mb-2">
 				Journal Entries
 			</h1>
 			<p class="text-gray-600">Capture your thoughts and memories</p>
 		</div>
 		<div class="w-full sm:w-auto">
 			<input
 				type="text"
 				placeholder="Search entries..."
 				bind:value={searchTerm}
 				class="w-full sm:w-64 px-3 py-2 border border-gray-300 rounded-md"
 			/>
 		</div>
 	</div>

 	{#if loading}
 		<div class="grid gap-6 md:grid-cols-2">
 			{#each Array(6) as _}
 				<div class="border rounded-lg p-6 bg-white shadow">
 					<div class="h-6 bg-gray-200 rounded mb-4 w-3/4"></div>
 					<div class="h-4 bg-gray-200 rounded mb-2 w-1/2"></div>
 					<div class="h-4 bg-gray-200 rounded mb-2"></div>
 					<div class="h-4 bg-gray-200 rounded mb-2"></div>
 					<div class="h-4 bg-gray-200 rounded w-3/4"></div>
 				</div>
 			{/each}
 		</div>
 	{:else if error}
 		<div class="bg-red-50 border border-red-200 rounded-lg p-4">
 			<p class="text-red-800">{error}</p>
 			<button class="mt-2 px-4 py-2 bg-red-100 text-red-800 rounded hover:bg-red-200" onclick={loadEntries}>
 				Try Again
 			</button>
 		</div>
 	{:else if filteredEntries.length === 0}
 		<div class="text-center py-16">
 			<h2 class="text-2xl font-bold text-gray-800 mb-3">
 				{searchTerm ? 'No matching entries found' : 'Your journal is empty'}
 			</h2>
 			<p class="text-gray-600 mb-8">
 				{searchTerm ? 'Try searching with different keywords.' : 'Start capturing your memories and thoughts.'}
 			</p>
 			<a href="/new" class="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700">Create Your First Entry</a>
 		</div>
 	{:else}
 		<div class="grid gap-6 md:grid-cols-2">
 			{#each filteredEntries as entry (entry.id)}
 				<div class="border rounded-lg p-6 bg-white shadow hover:shadow-lg transition-shadow duration-200">
 					<div class="flex justify-between items-start">
 						<div class="flex-1">
 							<h3 class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600 line-clamp-1">
 								<a href="/entry/{entry.id}" class="hover:from-indigo-600 hover:to-purple-600 transition-all duration-200">
 									{entry.title}
 								</a>
 							</h3>
 							<p class="text-sm text-gray-500 mt-2">
 								{formatDate(entry.createdAt)}
 							</p>
 						</div>
 					</div>
 					<p class="text-gray-700 line-clamp-3 mb-4 leading-relaxed mt-4">
 						{entry.content}
 					</p>
 					<div class="flex justify-between items-center">
 						<a href="/entry/{entry.id}" class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">View</a>
 						<div class="flex space-x-2">
 							<a href="/entry/{entry.id}/edit" class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">Edit</a>
 							<button onclick={() => deleteEntry(entry.id)} class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">
 								Delete
 							</button>
 						</div>
 					</div>
 				</div>
 			{/each}
 		</div>
 	{/if}
</div>