import redis
import threading


class RedisClient:
    _instances = {}
    _lock = threading.Lock()

    def __new__(cls, host="localhost", port=6379, db=0):
        instance_key = (host, port, db)
        if instance_key not in cls._instances:
            with cls._lock:
                if instance_key not in cls._instances:
                    instance = super(RedisClient, cls).__new__(cls)
                    instance.host = host
                    instance.port = port
                    instance.db = db
                    instance.redis_client = redis.StrictRedis(host=host, port=port, db=db)
                    cls._instances[instance_key] = instance
        return cls._instances[instance_key]

    def set(self, key, value):
        self.redis_client.set(key, value)

    def get(self, key):
        return self.redis_client.get(key)

    def delete(self, key):
        self.redis_client.delete(key)

    def keys(self, pattern="*"):
        return self.redis_client.keys(pattern)

    def ping(self):
        return self.redis_client.ping()
