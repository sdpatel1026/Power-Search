import os
import mysql.connector.pooling
from mysql.connector.errors import Error


class MySQLClient:
    _pool = None

    @classmethod
    def get_instance(cls):
        try:
            if cls._pool is None:
                cls._pool = mysql.connector.pooling.MySQLConnectionPool(
                    pool_name="my_pool",
                    pool_size=5,
                    host=os.getenv("MYSQL_HOST"),
                    user=os.getenv("MYSQL_USER"),
                    password=os.getenv("MYSQL_PASSWORD"),
                    database=os.getenv("MYSQL_DB")
                )
            conn = cls._pool.get_connection()
        except Error as e:
            print(e)  # log error --------------------
            return None
        finally:
            if conn and conn.is_connected():
                conn.close()
        return conn

    @staticmethod
    def close(connection, cursor):
        connection.close()
        cursor.close()

    @classmethod
    def execute(cls, sql, args=None, commit=False):
        """
            Execute a sql, it could be with args and without args. The usage is
            similar with execute() function in module pymysql.
            :param sql: sql clause
            :param args: args need by sql clause
            :param commit: whether to commit
            :return: if commit, return None, else, return result
            """
        # get connection form connection pool instead of create one.
        conn = cls.get_instance()
        cursor = conn.cursor()
        if args:
            cursor.execute(sql, args)
        else:
            cursor.execute(sql)
        if commit is True:
            conn.commit()
            MySQLClient.close(connection=conn, cursor=cursor)
            return None
        else:
            res = cursor.fetchall()
            MySQLClient.close(connection=conn, cursor=cursor)
            return res

    @classmethod
    def executemany(cls, sql, args, commit=False):
        """
            Execute with many args. Similar with executemany() function in pymysql.
            args should be a sequence.
            :param sql: sql clause
            :param args: args
            :param commit: commit or not.
            :return: if commit, return None, else, return result
            """
        # get connection form connection pool instead of create one.
        conn = cls.get_instance()
        cursor = conn.cursor()
        cursor.executemany(sql, args)
        if commit is True:
            conn.commit()
            MySQLClient.close(connection=conn, cursor=cursor)
            return None
        else:
            res = cursor.fetchall()
            MySQLClient.close(connection=conn, cursor=cursor)
            return res
