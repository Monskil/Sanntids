#include <stdio.h>
#include <pthread.h>



//pthread_mutex_t i_mut;
pthread_mutex_t mut = PTHREAD_MUTEX_INITIALIZER;
//int pthread_mutex_lock(pthread_mutex_t *mut)); 
//int pthread_mutex_unlock(pthread_mutex_t *mut);

int i = 0;
//int x;

void* someThreadFunction1(){
	int x;
	for (x = 0; x < 1000000; x++){
		pthread_mutex_lock(&mut);		
		i++;
		pthread_mutex_unlock(&mut);
	}
	
}

void* someThreadFunction2(){
	int x;
		
	for (x = 0; x < 1000002; x++){
		pthread_mutex_lock(&mut);
		i--;
		pthread_mutex_unlock(&mut);
	}
	
}

int main(){
	pthread_t thread_1;
	pthread_create(&thread_1, NULL, someThreadFunction1, NULL);

	pthread_t thread_2;
	pthread_create(&thread_2, NULL, someThreadFunction2, NULL);

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);

	printf("%d\n",i);
	return 0;
}