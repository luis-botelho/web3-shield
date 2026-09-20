import os
import psycopg2
from psycopg2.extras import RealDictCursor
from dotenv import load_dotenv

# Carrega a configuração local antes de abrir conexões com o banco.
load_dotenv()

class PostgresRepository:
    def __init__(self):
        self.conn = psycopg2.connect(os.getenv("DATABASE_URL"))
        
    def get_unprocessed_transactions(self):
        """
        Retorna transações sem análise de risco, indexadas por nome de coluna.
        """
        query = """
            SELECT r.tx_hash, r.function_signature 
            FROM raw_transactions r
            LEFT JOIN risk_analysis a ON r.tx_hash = a.tx_hash
            WHERE a.tx_hash IS NULL;
        """
        with self.conn.cursor(cursor_factory=RealDictCursor) as cursor:
            cursor.execute(query)
            return cursor.fetchall()
            
    def save_risk_analysis(self, tx_hash: str, risk_score: int):
        """
        Persiste a classificação sem sobrescrever análises existentes.
        """
        query = """
            INSERT INTO risk_analysis (tx_hash, risk_score)
            VALUES (%s, %s)
            ON CONFLICT (tx_hash) DO NOTHING;
        """
        with self.conn.cursor() as cursor:
            cursor.execute(query, (tx_hash, risk_score))
        self.conn.commit()

    def close(self):
        self.conn.close()
