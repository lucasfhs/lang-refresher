export namespace execution {
	
	export class Result {
	    stdout: string;
	    stderr: string;
	    exitCode: number;
	    durationMs: number;
	    timedOut: boolean;
	    cancelled: boolean;
	    passed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stdout = source["stdout"];
	        this.stderr = source["stderr"];
	        this.exitCode = source["exitCode"];
	        this.durationMs = source["durationMs"];
	        this.timedOut = source["timedOut"];
	        this.cancelled = source["cancelled"];
	        this.passed = source["passed"];
	    }
	}

}

export namespace main {
	
	export class ExerciseView {
	    id: string;
	    category: string;
	    title: string;
	    topic: string;
	    objective: string;
	    description: string;
	    requirements: string[];
	    example: string;
	    exampleOutput: string;
	    hint: string;
	    starterCode: string;
	    estimatedMinutes: number;
	    difficulty: number;
	    hasValidator: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExerciseView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.topic = source["topic"];
	        this.objective = source["objective"];
	        this.description = source["description"];
	        this.requirements = source["requirements"];
	        this.example = source["example"];
	        this.exampleOutput = source["exampleOutput"];
	        this.hint = source["hint"];
	        this.starterCode = source["starterCode"];
	        this.estimatedMinutes = source["estimatedMinutes"];
	        this.difficulty = source["difficulty"];
	        this.hasValidator = source["hasValidator"];
	    }
	}
	export class StateView {
	    trackId: string;
	    trackTitle: string;
	    trackDescription: string;
	    language: string;
	    currentIndex: number;
	    total: number;
	    completed: number;
	    skipped: number;
	    elapsedSeconds: number;
	    exercise: ExerciseView;
	    draft: string;
	    status: string;
	    statuses: Record<string, string>;
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new StateView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.trackId = source["trackId"];
	        this.trackTitle = source["trackTitle"];
	        this.trackDescription = source["trackDescription"];
	        this.language = source["language"];
	        this.currentIndex = source["currentIndex"];
	        this.total = source["total"];
	        this.completed = source["completed"];
	        this.skipped = source["skipped"];
	        this.elapsedSeconds = source["elapsedSeconds"];
	        this.exercise = this.convertValues(source["exercise"], ExerciseView);
	        this.draft = source["draft"];
	        this.status = source["status"];
	        this.statuses = source["statuses"];
	        this.locale = source["locale"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

