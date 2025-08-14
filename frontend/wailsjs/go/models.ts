export namespace service {
	
	export class Profiles {
	    id: number;
	    name: string;
	    hashtag: string;
	    first_comment: string;
	    is_authenticated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Profiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.hashtag = source["hashtag"];
	        this.first_comment = source["first_comment"];
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
	    ProfileDouyinID: number;
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
	        this.ProfileDouyinID = source["ProfileDouyinID"];
	        this.Status = source["Status"];
	    }
	}

}

