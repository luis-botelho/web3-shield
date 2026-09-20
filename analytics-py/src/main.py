import time
import sys

from src.repository.postgres_repo import PostgresRepository
from src.domain.risk_scorer import calculate_risk_score

def main():
    print("🕵️‍♂️ Iniciando o Detetive (Web3 Shield Analytics)...")
    
    try:
        repo = PostgresRepository()
        print("🐘 Conectado ao banco de dados com sucesso!")
    except Exception as e:
        print(f"🚨 Erro fatal ao conectar no banco: {e}")
        sys.exit(1)

    print("🔄 Vigiando o banco de dados em busca de novas transações...")

    try:
        while True:
            # Processa somente transações que ainda não possuem análise de risco.
            unprocessed = repo.get_unprocessed_transactions()
            
            if unprocessed:
                print(f"\n📦 Encontradas {len(unprocessed)} transações pendentes de análise.")
                
                for tx in unprocessed:
                    tx_hash = tx['tx_hash']
                    signature = tx['function_signature']
                    
                    score = calculate_risk_score(signature)
                    
                    repo.save_risk_analysis(tx_hash, score)
                    
                    print(f"  🔍 Analisado: {tx_hash[:10]}... | Sig: {signature} -> Nível de Risco: {score}")
            
            # Intervalo de polling do pipeline de análise.
            time.sleep(3)
            
    except KeyboardInterrupt:
        print("\n🛑 Encerrando o Detetive. Fechando conexão com o banco...")
    finally:
        repo.close()

if __name__ == "__main__":
    main()
