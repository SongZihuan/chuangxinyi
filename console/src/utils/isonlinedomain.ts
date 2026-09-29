export function IsOnlineDomain(url: string): boolean {
    return (new URL(url)).hostname == import.meta.env.VITE_ONLINE_DOMAIN
}