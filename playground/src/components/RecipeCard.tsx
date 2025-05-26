import React, { useCallback, useEffect, useState } from "react";
import { Button } from "./ui/button";
import { IconArrowDown, IconArrowUp, IconBrandInstagram, IconBrandTiktok, IconBrandYoutube, IconChevronDown, IconChevronRight, IconLoader2, IconMoodSadDizzy, IconRefresh, IconVideo, IconX } from "@tabler/icons-react";
import { toast } from "sonner";
import { playgroundConfig } from "@/config";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table";

export type RecipeType = 'tiktok' | 'instagram' | 'youtube' | 'generic';

const recipeTypeIcons: Record<RecipeType, React.ComponentType<{ size?: number }>> = {
	tiktok: IconBrandTiktok,
	instagram: IconBrandInstagram,
	youtube: IconBrandYoutube,
	generic: IconVideo,
};

const recipeTypeClasses: Record<RecipeType, string> = {
	tiktok: 'bg-black text-white',
	instagram: 'bg-gradient-to-br from-[#aa2cc0] via-[#ec0075] to-[#ffd366] text-white',
	youtube: 'bg-gradient-to-br from-[#fd0809] to-[#ff0033] text-white',
	generic: 'bg-gradient-to-br from-gray-100 to-gray-200 text-gray-800',
};

type ApiResponse = {
	recipe: {
		title: string;
		description: string;
		servings: number;
		prep_time: number;
		cook_time: number;
		total_time: number;
		difficulty: 1 | 2 | 3;
		ingredients: Array<{
			name: string;
			quantity: string;
			unit: string;
		}>;
		sections: Array<{
			instructions: Array<{
				optional: boolean;
				text: string;
			}>;
		}>;
		notes: string;
		nutritional_info: {
			calories: number;
			protein: number;
			carbohydrates: number;
			fats: number;
			fiber: number;
			sugar: number;
		};
		url: string;
	};
	metadata: {
		promptTokenCount: number;
		candidatesTokenCount: number;
	};
}


export function RecipeCard({ url, name, type, apiKey, index, setName }: { url: string; name: string; type: RecipeType; apiKey: string, index: number, setName: (name: string) => void }) {
	const [loading, setLoading] = useState(false);
	const [response, setResponse] = useState<ApiResponse | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [expanded, setExpanded] = useState(true);

	const fetchApi = useCallback(() => {
		if (!url) {
			return;
		}
		console.log('Fetching recipe for URL:', url);
		setLoading(true);
		setResponse(null);
		setError(null);

		fetch(`${playgroundConfig.apiRoot}/recipes/extract?url=${encodeURIComponent(url)}`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${apiKey}`,
			},
		})
			.then((response) => {
				if (response.ok) {
					return response.json();
				} else {
					throw new Error('Network response was not ok');
				}
			})
			.then((data) => {
				setName(data.recipe.title);
				setResponse(data);
				setError(null);
			})
			.catch((error) => {
				console.error('Error:', error);
				toast.error('Error: ' + error.message);
				setError(error.message);
			})
			.finally(() => {
				setLoading(false);
			});
	}, [apiKey, url, setName]);

	useEffect(() => {
		if (loading || error || response) {
			return;
		}
		fetchApi();
	}, [error, fetchApi, loading, response]);

	useEffect(() => {
		if(index === 0) {
			return;
		}
		setExpanded(false);
	}, [index]);

	console.log('loading', loading, 'error', error, 'response', response);
	const ChevronIcon = expanded ? IconChevronDown : IconChevronRight;
	return (
		<div className="rounded-lg border w-full p-4 space-y-4 shadow-lg bg-card">
			<div className="flex items-center">
				<div className="flex items-center gap-2">
					<Button variant='ghost' size='sm' className="px-1" onClick={() => setExpanded((prev) => !prev)}>
						<ChevronIcon className="text-primary" size={24} />
					</Button>
					{type && (
						<div className={`rounded-full p-2 ${recipeTypeClasses[type]} flex items-center justify-center`}>
							{React.createElement(recipeTypeIcons[type], { size: 24 })}
						</div>
					)}
					<h2 className="text-lg font-bold text-primary">{name}</h2>
				</div>
				{!expanded && !!error && (
					<div className="ml-auto rounded-full bg-gradient-to-br from-red-100 to-red-200 dark:from-red-700 dark:to-red-800 p-1">
						<IconX className="text-destructive" size={24} />
					</div>
				)}
				{!!response && (
					<div className="flex gap-2 justify-end ml-auto">
						<span className="bg-gradient-to-b from-gray-300 to-gray-400 text-gray-800 py-1 px-4 rounded-full text-sm font-medium">
							<span className="text-lg">{response.metadata.promptTokenCount}</span> <IconArrowUp className="inline" size={16} />
						</span>
						<span className="bg-gradient-to-b from-gray-300 to-gray-400 text-gray-800 py-1 px-4 rounded-full text-sm font-medium">
							<span className="text-lg">{response.metadata.candidatesTokenCount}</span> <IconArrowDown className="inline" size={16} />
						</span>
					</div>
				)}
			</div>
			{!!expanded && (
				<div className="grid md:grid-cols-2 gap-4 md:gap-8 p-4 min-h-[40vh]">
					{loading && (
						<IconLoader2 className="md:col-span-2 animate-spin my-auto mx-auto" size={32} />
					)}
					{!!response && (
						<>
							<div className="">
								<p className="font-semibold mb-2 text-lg">
									Description
								</p>
								<p className="text-foreground/70 mb-2">{response.recipe.description}</p>
							</div>
							<ul className="flex flex-wrap h-fit max-w-[300px]">
									<li className="flex w-full">
										<span className="text-foreground font-semibold mb-1">Servings</span>
										<div className="border-b border-dotted border-foreground/20 mb-2 mx-2 flex-1"></div>
										<span className="text-foreground/70">{response.recipe.servings}</span>
									</li>
									<li className="flex w-full">
										<span className="text-foreground font-semibold mb-1">Total Time</span>
										<div className="border-b border-dotted border-foreground/20 mb-2 mx-2 flex-1"></div>
										<span className="text-foreground/70">{response.recipe.total_time} m</span>
									</li>
									<li className="flex w-full">
										<span className="text-foreground font-light mb-1">Prep Time</span>
										<div className="border-b border-dotted border-foreground/20 mb-2 mx-2 flex-1"></div>
										<span className="text-foreground/70">{response.recipe.prep_time} m</span>
									</li>
									<li className="flex w-full">
										<span className="text-foreground font-light mb-1">Cook Time</span>
										<div className="border-b border-dotted border-foreground/20 mb-2 mx-2 flex-1"></div>
										<span className="text-foreground/70">{response.recipe.cook_time} m</span>
									</li>
							</ul>
							{/* <p className="mb-2"><span className="font-semibold">URL:</span> <a href={response.recipe.url} className="text-blue-600 underline" target="_blank" rel="noopener noreferrer">{response.recipe.url}</a></p> */}
							
							<div className="">
								<h4 className="font-semibold mb-2 text-lg">Ingredients</h4>
								<ul className="list-disc max-w-[300px]">
									{response.recipe.ingredients.map((ing, i) => (
										<li key={i} className="flex mb-2">
											<span className="">{ing.name}</span>
											<div className="border-b border-dotted border-foreground/20 mb-2 mx-2 flex-1"></div>
											<span className="text-foreground/70">{ing.quantity} {ing.unit}</span>
										</li>
									))}
								</ul>
							</div>
							<div/>
							<div className="md:col-span-2">
								<h4 className="font-semibold mb-2 text-lg">Instructions</h4>
								<ol className="list-decimal ml-6 divide-y divide-foreground/20 space-y-4">
									{response.recipe.sections.map((section, sectionIndex) => (
										<section key={sectionIndex} className="pb-4">
											{section.instructions.map((instruction, instructionIndex) => (
												<li key={instructionIndex} className="mb-2">
													{instruction.optional && <span className="text-yellow-500">[Optional] </span>}
													{instruction.text}
												</li>
											))}
										</section>
									))}
								</ol>
							</div>
							
							{!!response.recipe.notes && (
								<div className="">
									<h4 className="font-semibold mb-2 text-lg">Notes</h4>
									<p className="text-foreground/70">{response.recipe.notes}</p>
									<h4 className="font-semibold mb-2 text-lg mt-4">Source</h4>
									<a href={response.recipe.url} className="text-blue-600 underline" target="_blank" rel="noopener noreferrer">
										{response.recipe.url}
									</a>
								</div>
							)}
							<div className="">
								<h4 className="font-semibold mb-2 text-lg">Nutritional Information (100g)</h4>
								<div className="mt-4 bg-background overflow-hidden rounded-md border">
									<Table>
										<TableHeader>
											<TableRow>
												<TableHead className="bg-muted/50 py-2 font-medium">Nutrient</TableHead>
												<TableHead className="bg-muted/50 py-2 font-medium">Value</TableHead>
											</TableRow>
										</TableHeader>
										<TableBody>
											<TableRow>
												<TableCell className="py-2">Energy</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.calories} kcal</TableCell>
											</TableRow>
											<TableRow>
												<TableCell className="py-2">Protein</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.protein} g</TableCell>
											</TableRow>
											<TableRow>
												<TableCell className="py-2">Carbohydrates</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.carbohydrates} g</TableCell>
											</TableRow>
											<TableRow>
												<TableCell className="py-2">Fats</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.fats} g</TableCell>
											</TableRow>
											<TableRow>
												<TableCell className="py-2">Fiber</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.fiber} g</TableCell>
											</TableRow>
											<TableRow>
												<TableCell className="py-2">Sugar</TableCell>
												<TableCell className="py-2">{response.recipe.nutritional_info.sugar} g</TableCell>
											</TableRow>
										</TableBody>
									</Table>
								</div>
							</div>
						</>
					)}
					{!!error && (
						<div className="flex flex-col justify-center items-center gap-4 flex-1 md:col-span-2">
							<div className="p-5 rounded-full bg-gradient-to-br from-rose-100 to-red-100 dark:from-rose-700 dark:to-red-700">
								<div className="bg-gradient-to-br from-rose-500 to-red-500 dark:from-rose-500 dark:to-red-500 rounded-full p-4 shadow-2xl shadow-[inset_0_0_20px_rgba(255,255,255,0.5)]">
									<IconMoodSadDizzy
										stroke={2}
										className="text-success h-20 w-20 text-primary-foreground text-white opacity-85 [mask-image:linear-gradient(to_top,#00000060_0%,black_60%)] drop-shadow"
									/>
								</div>
							</div>
							<h1 className="text-2xl text-center font-bold text-primary">Error processing recipe</h1>
							<p className="text-muted-foreground text-center">
								{error}
								<br />
							Please try again or contact support.
							</p>
							<Button
								onClick={() => {
									setError(null);
									fetchApi();
								}}
								leftSection={<IconRefresh size={16} />}
								className="mt-4"
								variant="outline"
							>
							Retry
							</Button>
						</div>
					)}
				</div>
			)}
		</div>
	)
}