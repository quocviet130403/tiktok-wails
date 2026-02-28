export namespace service {
	
	export class ProfileDouyin {
	    id: number;
	    nickname: string;
	    url: string;
	    last_video_reup: any;
	    retry_count: number;
	    has_translate: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProfileDouyin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nickname = source["nickname"];
	        this.url = source["url"];
	        this.last_video_reup = source["last_video_reup"];
	        this.retry_count = source["retry_count"];
	        this.has_translate = source["has_translate"];
	    }
	}
	export class Profiles {
	    id: number;
	    name: string;
	    hashtag: string;
	    first_comment: string;
	    is_authenticated: boolean;
	    proxy_ip: string;
	    proxy_port: string;
	
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
	        this.proxy_ip = source["proxy_ip"];
	        this.proxy_port = source["proxy_port"];
	    }
	}
	export class Video {
	    id: number;
	    title: string;
	    video_url: string;
	    thumbnail_url: string;
	    duration: number;
	    like_count: number;
	    profile_douyin_id: number;
	    status: string;
	    is_deleted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.video_url = source["video_url"];
	        this.thumbnail_url = source["thumbnail_url"];
	        this.duration = source["duration"];
	        this.like_count = source["like_count"];
	        this.profile_douyin_id = source["profile_douyin_id"];
	        this.status = source["status"];
	        this.is_deleted = source["is_deleted"];
	    }
	}

}

