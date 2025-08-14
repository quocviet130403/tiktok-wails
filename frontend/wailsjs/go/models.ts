export namespace service {
	
	export class Accounts {
	    id: number;
	    name: string;
	    url_reup: string;
	    hashtag: string;
	    first_comment: string;
	    last_video_reup: string;
	    retry_count: number;
	    is_authenticated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Accounts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url_reup = source["url_reup"];
	        this.hashtag = source["hashtag"];
	        this.first_comment = source["first_comment"];
	        this.last_video_reup = source["last_video_reup"];
	        this.retry_count = source["retry_count"];
	        this.is_authenticated = source["is_authenticated"];
	    }
	}
	export class Video {
	    ID: number;
	    Title: string;
	    VideoURL: string;
	    ThumbnailURL: string;
	    Duration: number;
	    LikeCount: number;
	    AccountID: number;
	    Status: string;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.VideoURL = source["VideoURL"];
	        this.ThumbnailURL = source["ThumbnailURL"];
	        this.Duration = source["Duration"];
	        this.LikeCount = source["LikeCount"];
	        this.AccountID = source["AccountID"];
	        this.Status = source["Status"];
	    }
	}

}

