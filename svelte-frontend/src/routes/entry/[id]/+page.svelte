 <script lang="ts">
 	import { onMount } from 'svelte';
 	import { page } from '$app/stores';
 	import { journalApi } from '$lib/api/journal';
 	import type { JournalEntry } from '$lib/types/journal';

	let entry: JournalEntry | null = null;
	let loading = true;
	let error: string | null = null;
	let showDeleteDialog = false;

	onMount(async () => {
		await loadEntry();
	});

	async function loadEntry() {
		try {
			error = null;
			const id = parseInt($page.params.id || '0');
			entry = await journalApi.getEntry(id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load entry';
		} finally {
			loading = false;
		}
	}

	async function deleteEntry() {
		if (!entry) return;

		try {
			await journalApi.deleteEntry(entry.id);
			window.location.href = '/';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete entry';
		} finally {
			showDeleteDialog = false;
		}
	}

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="max-w-4xl mx-auto">
 	{#if loading}
 		<div class="bg-white border rounded-lg shadow p-6">
 			<div class="h-8 bg-gray-200 rounded mb-4 w-3/4"></div>
 			<div class="h-4 bg-gray-200 rounded mb-2 w-1/2"></div>
 			<div class="h-4 bg-gray-200 rounded mb-4"></div>
 			<div class="h-4 bg-gray-200 rounded mb-4"></div>
 			<div class="h-4 bg-gray-200 rounded mb-4"></div>
 			<div class="h-4 bg-gray-200 rounded w-3/4"></div>
 		</div>
 	{:else if error}
 		<div class="bg-red-50 border border-red-200 rounded-lg p-4">
 			<p class="text-red-800">{error}</p>
 			<a href="/" class="inline-block mt-2 px-4 py-2 bg-red-100 text-red-800 rounded hover:bg-red-200">Back to entries</a>
 		</div>
 	{:else if entry}
 		<div class="bg-white border rounded-lg shadow">
 			<div class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-t-lg p-6">
 				<div class="flex justify-between items-start">
 					<div class="flex-1">
 						<h1 class="text-3xl font-bold">{entry.title}</h1>
 						<div class="flex items-center gap-4 mt-3">
 							<span class="px-3 py-1 bg-blue-100 text-blue-800 rounded text-sm">
 								Created: {formatDate(entry.createdAt)}
 							</span>
 							{#if entry.updatedAt !== entry.createdAt}
 								<span class="px-3 py-1 bg-indigo-100 text-indigo-800 rounded text-sm">
 									Updated: {formatDate(entry.updatedAt)}
 								</span>
 							{/if}
 						</div>
 					</div>
 					<div class="flex space-x-3">
 						<a href="/entry/{entry.id}/edit" class="px-4 py-2 border border-white rounded hover:bg-white hover:text-blue-600 text-white">Edit</a>
 						<button onclick={() => showDeleteDialog = true} class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">Delete</button>
 					</div>
 				</div>
 			</div>
 			<div class="p-8">
 				<div class="prose prose-lg max-w-none">
 					<p class="whitespace-pre-wrap text-gray-800 leading-relaxed">{entry.content}</p>
 				</div>
 			</div>
 			<div class="bg-gray-50 px-8 py-6 border-t border-gray-200 rounded-b-lg">
 				<a href="/" class="inline-block px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">← Back to all entries</a>
 			</div>
 		</div>

 		{#if showDeleteDialog}
 			<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
 				<div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
 					<h2 class="text-xl font-bold mb-4">Delete Entry</h2>
 					<p class="text-gray-600 mb-6">Are you sure you want to delete "{entry.title}"? This action cannot be undone.</p>
 					<div class="flex justify-end space-x-3">
 						<button onclick={() => showDeleteDialog = false} class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">Cancel</button>
 						<button onclick={deleteEntry} class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">Delete</button>
 					</div>
 				</div>
 			</div>
 		{/if}
 	{/if}
 </div>