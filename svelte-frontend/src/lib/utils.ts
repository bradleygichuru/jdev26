import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import type { Component } from "svelte";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// Utility types for UI components
export type WithElementRef<T = HTMLElement> = T & {
	ref?: T | null;
	/** A Svelte action. */
	action?: (node: T) => {
		destroy?(): void;
	};
	/** A Svelte action. */
	use?: (node: T) => {
		destroy?(): void;
	};
};

export type WithoutChildren<T> = Omit<T, "children" | "child">;
export type WithoutChildrenOrChild<T> = Omit<T, "children" | "child">;

export type WithChild<T> = T & {
	children?: Component;
	child?: Component;
};

export type WithChildren<T> = T & {
	children?: Component;
};