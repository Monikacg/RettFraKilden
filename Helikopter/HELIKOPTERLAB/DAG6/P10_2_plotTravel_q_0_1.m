load('travel.mat')

figure('Name','Travel q = 0.1','NumberTitle','off')
plot(travel_measured(1,:), travel_measured(2,:))
title('Travel graph for weight q = 0.1')
xlabel('time [s]')
ylabel('travel [rad]')