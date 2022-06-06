import type { Mod } from '$lib/generated';
import { lockfileMods, manifestMods } from '$lib/store';
import { get } from 'svelte/store';

export interface OrderBy {
  name: string;
  func: (mod: PartialMod) => unknown,
}

export interface Filter {
  name: string;
  func: (mod: PartialMod) => boolean,
}

export const orderByOptions: OrderBy[] = [
  { name: 'Name', func: (mod: PartialMod) => mod.name.trim() },
  { name: 'Last updated', func: (mod: PartialMod) => Date.now() - Date.parse(mod.last_version_date) },
];

export const filterOptions: Filter[] = [
  { name: 'All mods', func: () => true },
  { name: 'Installed', func: (mod: PartialMod) => mod.mod_reference in get(manifestMods) },
  { name: 'Not installed', func: (mod: PartialMod) => !(mod.mod_reference in get(manifestMods)) },
  { name: 'Enabled', func: (mod: PartialMod) => mod.mod_reference in get(lockfileMods) },
  { name: 'Disabled', func: (mod: PartialMod) => !(mod.mod_reference in get(lockfileMods)) },
];

export type PartialMod = Pick<Mod, 'id' | 'mod_reference' | 'name' | 'logo' | 'short_description' | 'views' | 'downloads' | 'last_version_date'> & { authors: { user: Pick<User, 'username'> }[] };