import { useEffect, useRef, useState } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { IconComet, IconScan } from "@tabler/icons-react";
import { toast } from "sonner";
import { RecipeCard, RecipeType } from "./RecipeCard";

const LOCAL_STORAGE_APY_KEY = '__rvp_apiKey';

const getRecipePlaceholderNameAndType = (url: string): [string, RecipeType] => {
	const urlObj = new URL(url);
	const hostname = urlObj.hostname.toLowerCase();
	if (hostname.includes('tiktok.com')) {
		return ['TikTok Recipe', 'tiktok'];
	} else if (hostname.includes('instagram.com')) {
		return ['Instagram Recipe', 'instagram'];
	} else if (hostname.includes('youtube.com') || hostname.includes('youtu.be')) {
		return ['YouTube Recipe', 'youtube'];
	} else {
		return ['Recipe Video', 'generic'];
	}
};

export function Playground() {
	const apiKeyInput = useRef<HTMLInputElement>(null);
	const videoUrlInput = useRef<HTMLInputElement>(null);
	const [recipes, setRecipes] = useState<Array<{
		name: string;
		type: RecipeType;
		url: string;
	}>>([]);

	const scrollToTop = () => {
		window.scrollTo({
			top: 0,
			behavior: 'smooth',
		});
	}

	const onUpload = async () => {
		const apiKey = apiKeyInput.current?.value.trim();
		const videoUrl = videoUrlInput.current?.value.trim();

		if (!apiKey) {
			toast.error('Please enter your API key');
			return;
		}

		if (!videoUrl) {
			toast.error('Please enter a video URL');
			return;
		}

		const [name, type] = getRecipePlaceholderNameAndType(videoUrl);
		setRecipes((prev) => [{
			name: name,
			type: type,
			url: videoUrl,
		}, ...prev]);
		scrollToTop();
	}

	const persistApiKey = (key: string) => {
		localStorage.setItem(LOCAL_STORAGE_APY_KEY, key);
	}

	useEffect(() => {
		const apiKey = localStorage.getItem(LOCAL_STORAGE_APY_KEY);
		if (apiKeyInput.current) {
			apiKeyInput.current.value = apiKey || '';
		}
	}, []);


	return (
		<div className="flex flex-col items-center min-h-screen">
			<div className="flex gap-2 w-full sticky top-0 bg-background/80 backdrop-blur z-1 p-4 border-b">
				<Input placeholder="API key" className="flex-1" ref={apiKeyInput} type="password" onChange={(e) => {
					const input = e.target as HTMLInputElement;
					persistApiKey(input.value);
				}} />
				<Input
					placeholder="Url to video"
					className="flex-3"
					ref={videoUrlInput}
				/>
				<Button onClick={onUpload} leftSection={<IconScan size={16} />}>Parse</Button>
			</div>
			<div className="flex flex-col gap-4 mt-4 p-4 w-full flex-1">
				{recipes.map((recipe, i) => (
					<RecipeCard key={recipes.length - i} url={recipe.url} name={recipe.name} type={recipe.type} apiKey={apiKeyInput.current?.value || ''} index={i} setName={(name: string) => {
						setRecipes((prev) => {
							const newRecipes = [...prev];
							newRecipes[i] = {
								...newRecipes[i],
								name,
							};
							return newRecipes;
						})
					}} />
				))}
				{
					recipes.length === 0 && (
						<div className="flex flex-col justify-center items-center gap-4 flex-1">
							<div className="p-5 rounded-full bg-gradient-to-br from-teal-100 to-emerald-100 dark:from-teal-700 dark:to-emerald-700">
								<div className="bg-gradient-to-br from-teal-500 to-emerald-500 dark:from-teal-500 dark:to-emerald-500 rounded-full p-4 shadow-2xl shadow-[inset_0_0_20px_rgba(255,255,255,0.5)]">
									<IconComet
										stroke={2}
										className="text-success h-20 w-20 text-primary-foreground text-white opacity-85 [mask-image:linear-gradient(to_top,#00000060_0%,black_60%)] drop-shadow"
									/>
								</div>
							</div>
							<h1 className="text-2xl text-center font-bold text-primary">Welcome to Recipe Video Parser</h1>
							<p className="text-muted-foreground text-center">
								Parse any recipe from TikTok, Instagram, or YouTube videos.
							</p>

						</div>
					)
				}
			</div>
		</div>
	);
}