/*
 * Utility class for defining shared event handlers
 */

// Handler for main navigation panel open/close event
export const handleNavToggle = (e:any) => {
    const t =  e.target;
    t.expanded = !t.expanded;
}