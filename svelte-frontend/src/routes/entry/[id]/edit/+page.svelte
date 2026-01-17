 <script lang="ts">
 	import { onMount } from 'svelte';
 	import { page } from '$app/stores';
 	import { journalApi } from '$lib/api/journal';
 	import type { JournalEntry } from '$lib/types/journal';

	let entry: JournalEntry | null = null;
	let loading = true;
	let error: string | null = null;
	let saving = false;

	let title = '';
	let content = '';

	onMount(async () => {
		await loadEntry();
	});

	async function loadEntry() {
		try {
			error = null;
			const id = parseInt($page.params.id || '0');
			entry = await journalApi.getEntry(id);
			title = entry.title;
			content = entry.content;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load entry';
		} finally {
			loading = false;
		}
	}

	async function saveEntry() {
		if (!title.trim() || !content.trim()) {
			alert('Please fill in both title and content');
			return;
		}

		saving = true;
		try {
			const updatedEntry = await journalApi.updateEntry(entry!.id, {
				title: title.trim(),
				content: content.trim()
			});
			window.location.href = `/entry/${updatedEntry.id}`;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save entry';
		} finally {
			saving = false;
		}
	}

	function cancel() {
		if (confirm('Are you sure you want to discard your changes?')) {
			window.location.href = `/entry/${entry!.id}`;
		}
	}
</script>

<div class="max-w-4xl mx-auto">
 	{#if loading}
 		<div class="bg-white border rounded-lg shadow p-6">
 			<div class="h-8 bg-gray-200 rounded mb-4 w-3/4"></div>
 			<div class="space-y-6">
 				<div class="h-10 bg-gray-200 rounded"></div>
 				<div class="h-32 bg-gray-200 rounded"></div>
 			</div>
 		</div>
 	{:else if error}
 		<div class="bg-red-50 border border-red-200 rounded-lg p-4">
 			<p class="text-red-800">{error}</p>
 			<a href="/entry/{$page.params.id || ''}" class="inline-block mt-2 px-4 py-2 bg-red-100 text-red-800 rounded hover:bg-red-200">Back to entry</a>
 		</div>
 	{:else if entry}
 		<div class="bg-white border rounded-lg shadow">
 			<div class="bg-gradient-to-r from-orange-500 to-amber-600 text-white rounded-t-lg p-6">
 				<h1 class="text-3xl font-bold">Edit Entry</h1>
 				<p class="text-orange-100 mt-2">Update your thoughts and memories</p>
 			</div>
 			<div class="p-8">
 				{#if error}
 					<div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
 						<p class="text-red-800">{error}</p>
 					</div>
 				{/if}

 				<form onsubmit={saveEntry} class="space-y-6">
 					<div class="space-y-2">
 						<label for="title" class="text-lg font-semibold text-gray-800">Title</label>
 						<input
 							id="title"
 							type="text"
 							bind:value={title}
 							placeholder="Enter entry title..."
 							required
 							class="w-full px-3 py-2 border border-gray-300 rounded-md text-lg"
 						/>
 					</div>

 					<div class="space-y-2">
 						<label for="content" class="text-lg font-semibold text-gray-800">Content</label>
 						<textarea
 							id="content"
 							bind:value={content}
 							rows={12}
 							placeholder="Write your thoughts..."
 							required
 							class="w-full px-3 py-2 border border-gray-300 rounded-md text-lg leading-relaxed resize-none"
 						></textarea>
 					</div>

 					<div class="flex justify-between items-center pt-4">
 						<button type="button" onclick={cancel} class="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50">
 							Cancel
 						</button>
 						<div class="flex space-x-3">
 							<a href="/entry/{entry.id}" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">View Entry</a>
 							<button type="submit" disabled={saving} class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50">
 								{#if saving}
 									Saving...
 								{:else}
 									Save Changes
 								{/if}
 							</button>
 						</div>
 					</div>
 				</form>
 			</div>
 		</div>
  	{/if}
 </div>